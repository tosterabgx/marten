package client

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/tosterabgx/marten/internal/protocol"
)

var controlAddr = protocol.JoinAddr(protocol.DefaultServerAddr, protocol.ControlPort)

func performHandshake(enc *json.Encoder, dec *json.Decoder) (uint16, error) {
	var clientHello = protocol.ClientHello{DesiredPort: 0}
	var serverHello protocol.ServerHello

	if err := enc.Encode(clientHello); err != nil {
		return 0, err
	}

	if err := dec.Decode(&serverHello); err != nil {
		return 0, err
	}

	return serverHello.ActualPort, nil
}

func performConnectionAccept(dec *json.Decoder) (net.Conn, error) {
	newConnection := protocol.NewConnection{}
	if err := dec.Decode(&newConnection); err != nil {
		return nil, err
	}

	tunnelConn, err := net.Dial("tcp", controlAddr)
	if err != nil {
		return nil, err
	}

	acceptConnection := protocol.AcceptConnection{UUID: newConnection.UUID}

	tunnelEnc := json.NewEncoder(tunnelConn)
	if err := tunnelEnc.Encode(acceptConnection); err != nil {
		return nil, err
	}

	return tunnelConn, nil
}

func RunTCPTunnel(localPort uint16) error {
	controlConn, err := net.Dial("tcp", controlAddr)
	if err != nil {
		return err
	}
	defer controlConn.Close()

	enc := json.NewEncoder(controlConn)
	dec := json.NewDecoder(controlConn)

	actualPort, err := performHandshake(enc, dec)
	if err != nil {
		return err
	}

	localAddr := protocol.JoinAddr("localhost", localPort)
	actualAddr := protocol.JoinAddr(protocol.DefaultServerAddr, actualPort)
	fmt.Printf("Tunnel set:\n  %v -> %v\n", localAddr, actualAddr)

	for {
		tunnelConn, err := performConnectionAccept(dec)
		if err != nil {
			fmt.Printf("Error creating tunnel: %v\n", err)
			continue
		}

		localConn, err := net.Dial("tcp", localAddr)
		if err != nil {
			fmt.Printf("Error connecting to %v: %v\n", localAddr, err)
			continue
		}

		go protocol.Proxy(tunnelConn, localConn)
	}
}
