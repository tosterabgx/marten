package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/tosterabgx/marten/internal/protocol"
)

var controlAddr = protocol.JoinAddr(protocol.DefaultTunnelAddr, protocol.ControlPort)

func RunTunnel(localPort uint16, connType protocol.ConnType) error {
	controlConn, err := net.Dial("tcp", controlAddr)
	if err != nil {
		return errors.New("server unreachable")
	}
	defer controlConn.Close()

	enc := json.NewEncoder(controlConn)
	dec := json.NewDecoder(controlConn)

	clientHello := protocol.ClientHello{Type: connType}
	serverHello, err := performHandshake(enc, dec, clientHello)
	if err != nil {
		return err
	}

	localAddr := protocol.JoinAddr("localhost", localPort)

	var serverAddr string
	switch connType {
	case protocol.TypeHTTP:
		serverAddr = "https://" + serverHello.Subdomain + "." + protocol.DefaultServerAddr
	case protocol.TypeTCP:
		serverAddr = protocol.JoinAddr(protocol.DefaultServerAddr, serverHello.Port)
	}
	fmt.Printf("→ forwarding %v -> %v\n", serverAddr, localAddr)
	fmt.Println("→ tunnel established · ● live · press Ctrl+C to stop")

	for {
		tunnelConn, err := performConnectionAccept(dec)
		if err != nil {
			fmt.Printf("Error creating tunnel: %v\n", err)
			continue
		}

		localConn, err := net.Dial("tcp", localAddr)
		if err != nil {
			fmt.Printf("Error connecting to %v: %v\n", localAddr, err)
			tunnelConn.Close()
			continue
		}

		//fmt.Println("start proxy")

		go protocol.Proxy(tunnelConn, localConn)
	}
}
