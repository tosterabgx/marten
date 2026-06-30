package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/tosterabgx/marten/internal/protocol"
)

var controlAddr = protocol.JoinAddr(protocol.DefaultServerAddr, protocol.ControlPort)

func performHandshake(enc *json.Encoder, dec *json.Decoder) (uint16, error) {
	clientMessage, err := protocol.NewMessage(protocol.TypeClientHello, protocol.ClientHello{})
	if err != nil {
		return 0, err
	}

	var serverMessage protocol.Message

	if err := enc.Encode(clientMessage); err != nil {
		return 0, err
	}

	if err := dec.Decode(&serverMessage); err != nil {
		return 0, err
	}

	var serverHello protocol.ServerHello
	if err := serverMessage.Decode(&serverHello); err != nil {
		return 0, err
	}

	return *serverHello.Port, nil
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
		return errors.New("server unreachable")
	}
	defer controlConn.Close()

	enc := json.NewEncoder(controlConn)
	dec := json.NewDecoder(controlConn)

	serverPort, err := performHandshake(enc, dec)
	if err != nil {
		return err
	}

	localAddr := protocol.JoinAddr("localhost", localPort)
	serverAddr := protocol.JoinAddr(protocol.DefaultServerAddr, serverPort)
	fmt.Printf("Tunnel set:\n  %v -> %v\n", localAddr, serverAddr)

	for {
		_, err := performConnectionAccept(dec)
		if err != nil {
			fmt.Printf("Error creating tunnel: %v\n", err)
			continue
		}

		// localConn, err := net.Dial("tcp", localAddr)
		// if err != nil {
		// 	fmt.Printf("Error connecting to %v: %v\n", localAddr, err)
		// 	tunnelConn.Close()
		// 	continue
		// }

		// go protocol.Proxy(tunnelConn, localConn)
	}
}
