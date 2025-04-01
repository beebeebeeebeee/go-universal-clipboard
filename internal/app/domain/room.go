package domain

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

type Client struct {
	Room *Room
	Send chan []byte
}

type InfoPayload struct {
	ClientCount int `json:"client_count"`
}

type Room struct {
	ID          string
	clients     map[*Client]bool
	broadcast   chan []byte
	register    chan *Client
	unregister  chan *Client
	LastMessage []byte
	lastUpdated time.Time
	mu          sync.RWMutex
}

func NewRoom(id string) *Room {
	room := &Room{
		ID:          id,
		broadcast:   make(chan []byte),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		lastUpdated: time.Now(),
	}
	go room.run()
	return room
}

func (r *Room) broadcastClientCount() {
	clientCount := len(r.clients)
	payload := InfoPayload{ClientCount: clientCount}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling client count: %v", err)
		return
	}

	msg := NewMessage("info", string(payloadBytes))
	msgBytes, err := msg.ToJSON()
	if err != nil {
		log.Printf("Error creating info message: %v", err)
		return
	}

	for client := range r.clients {
		select {
		case client.Send <- msgBytes:
		default:
			close(client.Send)
			delete(r.clients, client)
		}
	}
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
			r.updateLastUpdated()
			r.broadcastClientCount()
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.Send)
				r.updateLastUpdated()
				r.broadcastClientCount()
			}
		case message := <-r.broadcast:
			// Parse the incoming message
			msg, err := FromJSON(message)
			if err != nil {
				log.Printf("Error parsing message in room %s: %v", r.ID, err)
				continue
			}

			// Only store and broadcast messages of type "message"
			if msg.Type == "message" {
				// Store the message payload
				r.LastMessage = message
				r.updateLastUpdated()

				for client := range r.clients {
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(r.clients, client)
					}
				}
			}
		}
	}
}

func (r *Room) AddClient(client *Client) {
	r.register <- client
}

func (r *Room) RemoveClient(client *Client) {
	r.unregister <- client
}

func (r *Room) BroadcastToOthers(sender *Client, message []byte) {
	r.broadcast <- message
}

func (r *Room) GetClientCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.clients)
}

func (r *Room) GetLastUpdated() time.Time {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.lastUpdated
}

func (r *Room) updateLastUpdated() {
	r.mu.Lock()
	r.lastUpdated = time.Now()
	r.mu.Unlock()
}
