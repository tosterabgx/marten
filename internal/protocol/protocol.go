package protocol

import "github.com/google/uuid"

const ControlPort uint16 = 6472

// const DefaultServerAddr = "marten.tosterabgx.me"
const DefaultServerAddr = "localhost"

type ClientHello struct {
	DesiredPort uint16 `json:"ClientHello"`
}

type ServerHello struct {
	ActualPort uint16 `json:"ServerHello"`
}

type NewConnection struct {
	UUID uuid.UUID `json:"NewConnection"`
}

type AcceptConnection struct {
	UUID uuid.UUID `json:"AcceptConnection"`
}
