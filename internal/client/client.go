package client

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/tosterabgx/marten/internal/protocol"
)

func RunTCPTunnel(port uint16) error {
	conn, err := net.Dial("tcp", net.JoinHostPort(protocol.DefaultServerAddr, strconv.Itoa(int(protocol.ControlPort))))
	if err != nil {
		return err
	}
	defer conn.Close()

	var clientHello = protocol.ClientHello{RequiredPort: port}
	var serverHello protocol.ServerHello

	enc := json.NewEncoder(conn)
	if err := enc.Encode(clientHello); err != nil {
		return err
	}

	fmt.Printf("Sent ClientHello: %v\n", clientHello)

	dec := json.NewDecoder(conn)
	if err := dec.Decode(&serverHello); err != nil {
		return err
	}

	fmt.Printf("Got ServerHello: %v\n", serverHello)

	return nil
}
