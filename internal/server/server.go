package server

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/tosterabgx/marten/internal/protocol"
)

func handleNewClient(conn net.Conn, clientHello protocol.ClientHello) {
	fmt.Println("Got ClientHello:", clientHello)

	actualPort, _ := registerListener(clientHello.DesiredPort, conn) // TODO: handle error

	serverHello := protocol.ServerHello{ActualPort: actualPort}

	enc := json.NewEncoder(conn)
	if err := enc.Encode(serverHello); err != nil {
		fmt.Println("Error encoding ServerHello:", err)
		conn.Close()
		return
	}

	fmt.Printf("Sent ServerHello: %v\n", serverHello)
}

func handleAcceptConnection(tunnelConn net.Conn, acceptConnection protocol.AcceptConnection) {
	fmt.Println("Got AcceptConnection:", acceptConnection)
	uuid := acceptConnection.UUID

	connsMu.RLock()
	outboundConn, ok := outboundConnections[uuid]
	connsMu.RUnlock()
	if !ok {
		fmt.Println("No outbound connection found for", uuid)
		return
	}

	protocol.Proxy(tunnelConn, outboundConn)
}

func handleConnection(conn net.Conn) {
	dec := json.NewDecoder(conn)

	var raw map[string]json.RawMessage
	if err := dec.Decode(&raw); err != nil {
		fmt.Println("Incorrect message:", err)
		return
	}

	if data, ok := raw["ClientHello"]; ok {
		var clientHello protocol.ClientHello
		if err := json.Unmarshal(data, &clientHello.DesiredPort); err != nil {
			fmt.Println("Incorrect ClientHello:", err)
			conn.Close()
			return
		}

		handleNewClient(conn, clientHello)
		return
	}

	if data, ok := raw["AcceptConnection"]; ok {
		var acceptConnection protocol.AcceptConnection
		if err := json.Unmarshal(data, &acceptConnection.UUID); err != nil {
			fmt.Println("Incorrect AcceptConnection:", err)
			conn.Close()
			return
		}

		handleAcceptConnection(conn, acceptConnection)
		return
	}

	fmt.Println("Unknown message:", raw)
}

func RunControlServer() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", protocol.ControlPort))
	if err != nil {
		return err
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		go handleConnection(conn)
	}
}
