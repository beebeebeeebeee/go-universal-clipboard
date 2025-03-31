package domain

import (
	"sync"
	"time"
)

type Room struct {
	ID          string
	Clients     map[*Client]bool
	LastMessage []byte
	LastUpdated time.Time
	mu          sync.RWMutex
}

type Client struct {
	Room *Room
	Send chan []byte
}

func NewRoom(id string) *Room {
	return &Room{
		ID:          id,
		Clients:     make(map[*Client]bool),
		LastMessage: []byte{},
		LastUpdated: time.Now(),
	}
}

func (r *Room) Broadcast(message []byte) {
	r.mu.Lock()
	r.LastMessage = message
	r.LastUpdated = time.Now()
	r.mu.Unlock()

	r.mu.RLock()
	defer r.mu.RUnlock()

	for client := range r.Clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(r.Clients, client)
		}
	}
}

func (r *Room) BroadcastToOthers(sender *Client, message []byte) {
	r.mu.Lock()
	r.LastMessage = message
	r.LastUpdated = time.Now()
	r.mu.Unlock()

	r.mu.RLock()
	defer r.mu.RUnlock()

	for client := range r.Clients {
		if client != sender {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(r.Clients, client)
			}
		}
	}
}

func (r *Room) AddClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Clients[client] = true
	r.LastUpdated = time.Now()
}

func (r *Room) RemoveClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.Clients[client]; ok {
		delete(r.Clients, client)
		r.LastUpdated = time.Now()
	}
	close(client.Send)
}

// GetClientCount returns the number of clients in the room
func (r *Room) GetClientCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.Clients)
}

// GetLastUpdated returns the last update time of the room
func (r *Room) GetLastUpdated() time.Time {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.LastUpdated
}
