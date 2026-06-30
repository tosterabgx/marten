package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"

	"github.com/tosterabgx/marten/internal/protocol"
)

func handleNewClient(conn net.Conn, clientHello protocol.ClientHello) (net.Listener, error) {
	slog.Debug("got ClientHello", "message", clientHello)

	l, port, err := registerListener(conn)
	if err != nil {
		return nil, fmt.Errorf("listener registry failed: %v", err)
	}

	serverMessage, err := protocol.NewMessage(protocol.TypeServerHello, protocol.ServerHello{Port: &port})
	if err != nil {
		return nil, nil
	}

	enc := json.NewEncoder(conn)
	if err := enc.Encode(serverMessage); err != nil {
		return nil, fmt.Errorf("ServerHello encoding failed: %v", err)
	}

	slog.Debug("sent ServerHello", "message", serverMessage)
	return l, nil
}

func handleAcceptConnection(tunnelConn net.Conn, acceptConnection protocol.AcceptConnection) error {
	slog.Debug("got AcceptConnection", "message", acceptConnection)

	uuid := acceptConnection.UUID
	connsMu.RLock()
	externalConn, ok := externalConns[uuid]
	connsMu.RUnlock()
	if !ok {
		return fmt.Errorf("no external connection found with uuid=%v", uuid)
	}

	slog.Debug("start proxy")

	protocol.Proxy(tunnelConn, externalConn)
	connsMu.Lock()
	delete(externalConns, uuid)
	connsMu.Unlock()
	return nil
}
