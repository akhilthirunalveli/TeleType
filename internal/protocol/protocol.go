package protocol

import "time"

// MessageType defines the type of message being sent
type MessageType string

const (
	MsgTypeChat   MessageType = "CHAT"
	MsgTypeSystem MessageType = "SYSTEM"
	MsgTypeJoin   MessageType = "JOIN"
	MsgTypeLeave  MessageType = "LEAVE"
	MsgTypeName   MessageType = "NAME"
	MsgTypeError  MessageType = "ERROR"
)

// Message represents the data structure exchanged over WebSocket
type Message struct {
	Type      MessageType `json:"type"`
	Content   string      `json:"content,omitempty"`
	Sender    string      `json:"sender,omitempty"`
	Room      string      `json:"room,omitempty"`
	Timestamp time.Time   `json:"timestamp,omitempty"`
}

// NewChatMessage creates a standard chat message
func NewChatMessage(sender, content, room string) Message {
	return Message{
		Type:      MsgTypeChat,
		Content:   content,
		Sender:    sender,
		Room:      room,
		Timestamp: time.Now(),
	}
}

// NewSystemMessage creates a system notification
func NewSystemMessage(content string) Message {
	return Message{
		Type:      MsgTypeSystem,
		Content:   content,
		Sender:    "SYSTEM",
		Timestamp: time.Now(),
	}
}
