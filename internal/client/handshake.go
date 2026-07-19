package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"

	"github.com/tosterabgx/marten/internal/protocol"
)

func performHandshake(enc *json.Encoder, dec *json.Decoder, clientHello protocol.ClientHello) (protocol.ServerHello, error) {
	clientMessage := protocol.NewMessage(clientHello)

	var serverMessage protocol.Message

	if err := enc.Encode(clientMessage); err != nil {
		return protocol.ServerHello{}, err
	}

	if err := dec.Decode(&serverMessage); err != nil {
		return protocol.ServerHello{}, err
	}

	serverHello, ok := serverMessage.Payload.(protocol.ServerHello)
	if !ok {
		return protocol.ServerHello{}, fmt.Errorf("expected ServerHello, but received unexpected message type: %T", serverMessage.Payload)
	}

	return serverHello, nil
}

func performConnectionAccept(dec *json.Decoder) (net.Conn, error) {
	var msg protocol.Message
	if err := dec.Decode(&msg); err != nil {
		return nil, err
	}

	newConnection, ok := msg.Payload.(protocol.NewConnection)
	if !ok {
		return nil, fmt.Errorf("expected NewConnection, but received unexpected message type: %T", msg.Payload)
	}

	tunnelConn, err := tls.Dial("tcp", controlAddr, &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         protocol.DefaultTunnelAddr,
	})
	if err != nil {
		return nil, err
	}

	acceptConnection := protocol.NewMessage(protocol.AcceptConnection{UUID: newConnection.UUID})

	if err := json.NewEncoder(tunnelConn).Encode(acceptConnection); err != nil {
		return nil, err
	}

	return tunnelConn, nil
}
