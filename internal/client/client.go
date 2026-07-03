package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/tosterabgx/marten/internal/protocol"
)

var controlAddr = protocol.JoinAddr(protocol.DefaultServerAddr, protocol.ControlPort)

func RunTunnel(localPort uint16, isHttp bool) error {
	controlConn, err := net.Dial("tcp", controlAddr)
	if err != nil {
		return errors.New("server unreachable")
	}
	defer controlConn.Close()

	enc := json.NewEncoder(controlConn)
	dec := json.NewDecoder(controlConn)

	clientHelloType := "tcp"
	if isHttp {
		clientHelloType = "http"
	}

	clientHello := protocol.ClientHello{Type: clientHelloType}
	serverHello, err := performHandshake(enc, dec, clientHello)
	if err != nil {
		return err
	}

	localAddr := protocol.JoinAddr("localhost", localPort)

	var serverAddr string
	if isHttp {
		serverAddr = "http://" + serverHello.Subdomain + "." + protocol.DefaultServerAddr
	} else {
		serverAddr = protocol.JoinAddr(protocol.DefaultServerAddr, serverHello.Port)
	}
	fmt.Printf("Tunnel set:\n  %v -> %v\n", localAddr, serverAddr)

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

		fmt.Println("start proxy")

		go protocol.Proxy(tunnelConn, localConn)
	}
}
