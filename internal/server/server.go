package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"

	"github.com/tosterabgx/marten/internal/protocol"
)

func handleNewClient(conn net.Conn, clientHello protocol.ClientHello) (net.Listener, error) {
	slog.Debug("got ClientHello", "message", clientHello)

	l, actualPort, err := registerListener(conn)
	if err != nil {
		return nil, fmt.Errorf("listener registry failed: %v", err)
	}
	serverHello := protocol.ServerHello{ActualPort: actualPort}

	enc := json.NewEncoder(conn)
	if err := enc.Encode(serverHello); err != nil {
		return nil, fmt.Errorf("ServerHello encoding failed: %v", err)
	}

	slog.Debug("sent ServerHello", "message", serverHello)
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

func handleConnection(conn net.Conn) {
	dec := json.NewDecoder(conn)

	var raw map[string]json.RawMessage
	if err := dec.Decode(&raw); err != nil {
		slog.Warn("incorrect message", "error", err)
		conn.Close()
		return
	}

	if data, ok := raw["ClientHello"]; ok {
		var clientHello protocol.ClientHello
		if err := json.Unmarshal(data, &clientHello.DesiredPort); err != nil {
			slog.Warn("incorrect ClientHello", "error", err)
			conn.Close()
			return
		}

		l, err := handleNewClient(conn, clientHello)
		if err != nil {
			slog.Warn("failed to handle ClientHello", "error", err)
			conn.Close()
		}

		io.Copy(io.Discard, conn)
		l.Close()
	} else if data, ok := raw["AcceptConnection"]; ok {
		var acceptConnection protocol.AcceptConnection
		if err := json.Unmarshal(data, &acceptConnection.UUID); err != nil {
			slog.Warn("incorrect AcceptConnection", "error", err)
			conn.Close()
			return
		}

		err := handleAcceptConnection(conn, acceptConnection)
		if err != nil {
			slog.Warn("failed to handle AcceptConnection", "error", err)
			conn.Close()
		}
	} else {
		slog.Warn("unknown message", "message", raw)
	}
}

func RunControlServer() error {
	addr := protocol.JoinAddr("", protocol.ControlPort)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer l.Close()

	slog.Info("listening control", "address", addr)

	for {
		conn, err := l.Accept()
		if err != nil {
			slog.Error("error accepting", "error", err)
			continue
		}

		go handleConnection(conn)
	}
}
