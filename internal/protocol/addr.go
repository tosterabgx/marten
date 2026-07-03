package protocol

import (
	"net"
	"strconv"
)

const HTTPPort uint16 = 8080
const ControlPort uint16 = 6472
const MinPort uint16 = 6000
const MaxPort uint16 = 8000

const DefaultServerAddr = "marten.tosterabgx.me"

// const DefaultServerAddr = "localhost"

func JoinAddr(addr string, port uint16) string {
	addr = net.JoinHostPort(addr, strconv.Itoa(int(port)))
	return addr
}
