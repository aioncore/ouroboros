package pubsub

import "fmt"

type Message interface {
	GetType() MessageType
	String() string
}

type MessageType string

const (
	TX MessageType = "tx"
)

type MessageHeader struct {
	ID   string      `json:"id"`
	Type MessageType `json:"type"`
}

func (m *MessageHeader) GetType() MessageType {
	return m.Type
}

func (m *MessageHeader) GetID() string {
	return m.ID
}

func (m *MessageHeader) String() string {
	return fmt.Sprintf("%s message %s", m.Type, m.ID)
}
