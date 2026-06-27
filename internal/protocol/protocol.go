package protocol

import (
	"net"
	"strconv"

	"github.com/google/uuid"
)

const ControlPort uint16 = 6472
const MinPort uint16 = 6000
const MaxPort uint16 = 8000

const DefaultServerAddr = "marten.tosterabgx.me"

func JoinAddr(addr string, port uint16) string {
	addr = net.JoinHostPort(addr, strconv.Itoa(int(port)))
	return addr
}

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
