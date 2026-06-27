package server

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/tosterabgx/marten/internal/protocol"
)

func handleControlConnection(conn net.Conn) {
	var clientHello protocol.ClientHello

	dec := json.NewDecoder(conn)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&clientHello); err != nil {
		fmt.Println("not a ClientHello:", err)
		return
	}

	fmt.Printf("Got ClientHello: %v\n", clientHello)
	actualPort, _ := registerListener(clientHello.DesiredPort, conn) // TODO: handle error

	serverHello := protocol.ServerHello{ActualPort: actualPort}

	enc := json.NewEncoder(conn)
	if err := enc.Encode(serverHello); err != nil {
		fmt.Println("couldn't encode a ServerHello:", err)
		return
	}

	fmt.Printf("Sent ServerHello: %v\n", serverHello)
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

		go handleControlConnection(conn)
	}
}
