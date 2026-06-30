package protocol

import (
	"encoding/json"

	"github.com/google/uuid"
)

type MessageType string

const (
	TypeClientHello      MessageType = "ClientHello"
	TypeServerHello      MessageType = "ServerHello"
	TypeNewConnection    MessageType = "NewConnection"
	TypeAcceptConnection MessageType = "AcceptConnection"
)

type anyMessage interface {
	ClientHello | ServerHello | NewConnection | AcceptConnection
}

type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func (m Message) Decode(v any) error {
	return json.Unmarshal(m.Payload, v)
}

func NewMessage[T anyMessage](t MessageType, payload T) (Message, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return Message{}, err
	}
	return Message{Type: t, Payload: data}, nil
}

type ClientHello struct {
}

type ServerHello struct {
	Port      *uint16 `json:"port,omitempty"`
	Subdomain *string `json:"subdomain,omitempty"`
}

type NewConnection struct {
	UUID uuid.UUID `json:"uuid"`
}

type AcceptConnection struct {
	UUID uuid.UUID `json:"uuid"`
}
