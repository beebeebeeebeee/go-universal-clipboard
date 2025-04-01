package domain

import "encoding/json"

// Message represents the structure of WebSocket messages
type Message struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

// NewMessage creates a new message with the given type and payload
func NewMessage(msgType string, payload string) *Message {
	return &Message{
		Type:    msgType,
		Payload: payload,
	}
}

// ToJSON converts the message to JSON bytes
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// FromJSON parses JSON bytes into a Message
func FromJSON(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}
