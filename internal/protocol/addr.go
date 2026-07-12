package protocol

import (
	"net"
	"strconv"
)

const APIPort uint16 = 8081
const HTTPPort uint16 = 8080
const ControlPort uint16 = 6472
const MinPort uint16 = 10000
const MaxPort uint16 = 12000

const DefaultServerAddr = "usemarten.tech"

func JoinAddr(addr string, port uint16) string {
	addr = net.JoinHostPort(addr, strconv.Itoa(int(port)))
	return addr
}
