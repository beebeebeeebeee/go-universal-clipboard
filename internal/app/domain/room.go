package domain

import (
	"sync"
)

type Room struct {
	ID      string
	Clients map[*Client]bool
	mu      sync.RWMutex
}

type Client struct {
	Room *Room
	Send chan []byte
}

func NewRoom(id string) *Room {
	return &Room{
		ID:      id,
		Clients: make(map[*Client]bool),
	}
}

func (r *Room) Broadcast(message []byte) {
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
}

func (r *Room) RemoveClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.Clients[client]; ok {
		delete(r.Clients, client)
	}
	close(client.Send)
}
