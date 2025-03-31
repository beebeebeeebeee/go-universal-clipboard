package controller

import (
	"go-universal-clipboard/internal/app/domain"
	"go-universal-clipboard/internal/cfg"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// RoomInfo represents the information about a room for the API response
type RoomInfo struct {
	ID          string    `json:"id"`
	ClientCount int       `json:"client_count"`
	LastMessage string    `json:"last_message"`
	LastUpdated time.Time `json:"last_updated"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketController struct {
	rooms               map[string]*domain.Room
	mu                  sync.RWMutex
	emptyCleanupMinutes int
	idleCleanupMinutes  int
	cleanupTickerDone   chan bool
}

func NewWebSocketController() *WebSocketController {
	controller := &WebSocketController{
		rooms:               make(map[string]*domain.Room),
		emptyCleanupMinutes: cfg.Cfg.App.RoomEmptyCleanupMinutes,
		idleCleanupMinutes:  cfg.Cfg.App.RoomIdleCleanupMinutes,
		cleanupTickerDone:   make(chan bool),
	}

	// Start the cleanup ticker
	go controller.startCleanupTicker()

	return controller
}

func (c *WebSocketController) Setup(g *gin.RouterGroup) {
	rg := g.Group("")
	{
		rg.GET("", c.HandleWebSocket)
		rg.GET("/info", c.HandleRoomInfo)
	}
}

func (c *WebSocketController) startCleanupTicker() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanupInactiveRooms()
		case <-c.cleanupTickerDone:
			return
		}
	}
}

func (c *WebSocketController) cleanupInactiveRooms() {
	c.mu.Lock()
	defer c.mu.Unlock()

	emptyThreshold := time.Now().Add(-time.Duration(c.emptyCleanupMinutes) * time.Minute)
	idleThreshold := time.Now().Add(-time.Duration(c.idleCleanupMinutes) * time.Minute)

	for id, room := range c.rooms {
		if room.GetLastUpdated().Before(idleThreshold) {
			log.Printf("Removing idle room: %s (last updated: %s)", id, room.GetLastUpdated().Format(time.RFC3339))
			delete(c.rooms, id)
			continue
		}

		if room.GetClientCount() == 0 && room.GetLastUpdated().Before(emptyThreshold) {
			log.Printf("Removing inactive room: %s (last updated: %s)", id, room.GetLastUpdated().Format(time.RFC3339))
			delete(c.rooms, id)
		}
	}
}

func (c *WebSocketController) getOrCreateRoom(roomID string) *domain.Room {
	c.mu.Lock()
	defer c.mu.Unlock()

	if room, exists := c.rooms[roomID]; exists {
		return room
	}

	room := domain.NewRoom(roomID)
	c.rooms[roomID] = room
	log.Printf("Created new room: %s", roomID)
	return room
}

func (c *WebSocketController) HandleWebSocket(ctx *gin.Context) {
	roomID := ctx.Query("room")
	if roomID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "room parameter is required"})
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}

	room := c.getOrCreateRoom(roomID)
	client := &domain.Client{
		Room: room,
		Send: make(chan []byte, 256),
	}

	client.Room.AddClient(client)

	// Send the latest message to the new client if it exists
	if len(room.LastMessage) > 0 {
		select {
		case client.Send <- room.LastMessage:
		default:
			// If the send channel is full, we'll skip sending the latest message
			log.Printf("Failed to send latest message to new client in room %s: send buffer full", roomID)
		}
	}

	go c.writePump(client, conn)
	go c.readPump(client, conn)
}

func (c *WebSocketController) readPump(client *domain.Client, conn *websocket.Conn) {
	defer func() {
		client.Room.RemoveClient(client)
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// log error
			}
			break
		}

		client.Room.BroadcastToOthers(client, message)
	}
}

func (c *WebSocketController) writePump(client *domain.Client, conn *websocket.Conn) {
	defer func() {
		conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}

// HandleRoomInfo returns information about all active rooms
func (c *WebSocketController) HandleRoomInfo(ctx *gin.Context) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	roomInfos := make([]RoomInfo, 0, len(c.rooms))
	for id, room := range c.rooms {
		roomInfos = append(roomInfos, RoomInfo{
			ID:          id,
			ClientCount: room.GetClientCount(),
			LastMessage: string(room.LastMessage),
			LastUpdated: room.GetLastUpdated(),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"rooms": roomInfos,
		"total": len(roomInfos),
	})
}
