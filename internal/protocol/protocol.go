package protocol

import (
	"encoding/json"
	"fmt"

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
	getType() MessageType
}

type Message struct {
	Type    MessageType `json:"type"`
	Payload anyMessage  `json:"payload"`
}

func (m *Message) UnmarshalJSON(data []byte) error {
	var shadow struct {
		Type    MessageType     `json:"type"`
		Payload json.RawMessage `json:"payload"`
	}

	if err := json.Unmarshal(data, &shadow); err != nil {
		return err
	}

	m.Type = shadow.Type

	switch shadow.Type {
	case TypeClientHello:
		var p ClientHello
		if err := json.Unmarshal(shadow.Payload, &p); err != nil {
			return err
		}
		m.Payload = p

	case TypeServerHello:
		var p ServerHello
		if err := json.Unmarshal(shadow.Payload, &p); err != nil {
			return err
		}
		m.Payload = p

	case TypeNewConnection:
		var p NewConnection
		if err := json.Unmarshal(shadow.Payload, &p); err != nil {
			return err
		}
		m.Payload = p

	case TypeAcceptConnection:
		var p AcceptConnection
		if err := json.Unmarshal(shadow.Payload, &p); err != nil {
			return err
		}
		m.Payload = p

	default:
		return fmt.Errorf("unknown message type: %v", shadow.Type)
	}

	return nil
}

func NewMessage(payload anyMessage) Message {
	return Message{
		Type:    payload.getType(),
		Payload: payload,
	}
}

type ClientHello struct {
	Type string `json:"type"`
}

type ServerHello struct {
	Port      uint16 `json:"port,omitempty"`
	Subdomain string `json:"subdomain,omitempty"`
}

type NewConnection struct {
	UUID uuid.UUID `json:"uuid"`
}

type AcceptConnection struct {
	UUID uuid.UUID `json:"uuid"`
}

func (ClientHello) getType() MessageType      { return TypeClientHello }
func (ServerHello) getType() MessageType      { return TypeServerHello }
func (NewConnection) getType() MessageType    { return TypeNewConnection }
func (AcceptConnection) getType() MessageType { return TypeAcceptConnection }
