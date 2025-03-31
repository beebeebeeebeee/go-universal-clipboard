package controller

import (
	"go-universal-clipboard/internal/app/domain"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketController struct {
	rooms map[string]*domain.Room
	mu    sync.RWMutex
}

func NewWebSocketController() *WebSocketController {
	return &WebSocketController{
		rooms: make(map[string]*domain.Room),
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

	client := &domain.Client{
		Room: c.getOrCreateRoom(roomID),
		Send: make(chan []byte, 256),
	}

	client.Room.AddClient(client)

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
