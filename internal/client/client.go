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

	var clientHello = protocol.ClientHello{DesiredPort: port}
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

	for {
		var newConnection = protocol.NewConnection{}
		if err := dec.Decode(&newConnection); err != nil {
			return err
		}

		fmt.Printf("Got NewConnection: %v\n", newConnection)

		conn2, err := net.Dial("tcp", net.JoinHostPort(protocol.DefaultServerAddr, strconv.Itoa(int(protocol.ControlPort))))
		if err != nil {
			return err
		}

		enc2 := json.NewEncoder(conn2)

		var acceptConnection = protocol.AcceptConnection{UUID: newConnection.UUID}
		if err := enc2.Encode(acceptConnection); err != nil {
			return err
		}

		fmt.Printf("Sent AcceptConnection: %v\n", acceptConnection)

		localConn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			return fmt.Errorf("ERROR")
		}

		protocol.Proxy(conn2, localConn)

		fmt.Println("Connection ended")
	}
	return nil
}
