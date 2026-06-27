package server

import (
	"encoding/json"
	"log/slog"
	"net"

	"github.com/tosterabgx/marten/internal/protocol"
)

func handleNewClient(conn net.Conn, clientHello protocol.ClientHello) {
	slog.Debug("got ClientHello", "message", clientHello)

	actualPort, err := registerListener(conn)
	if err != nil {
		slog.Error("listener registry failed", "error", err)
		return
	}
	serverHello := protocol.ServerHello{ActualPort: actualPort}

	enc := json.NewEncoder(conn)
	if err := enc.Encode(serverHello); err != nil {
		slog.Warn("ServerHello encoding failed", "error", err)
		conn.Close()
		return
	}

	slog.Debug("sent ServerHello", "message", serverHello)
}

func handleAcceptConnection(tunnelConn net.Conn, acceptConnection protocol.AcceptConnection) {
	slog.Debug("got AcceptConnection", "message", acceptConnection)

	uuid := acceptConnection.UUID
	connsMu.RLock()
	externalConn, ok := externalConns[uuid]
	connsMu.RUnlock()
	if !ok {
		slog.Warn("no external connection found", "UUID", uuid)
		return
	}

	slog.Debug("start proxy")

	protocol.Proxy(tunnelConn, externalConn)
	connsMu.Lock()
	delete(externalConns, uuid)
	connsMu.Unlock()
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

		handleNewClient(conn, clientHello)
	} else if data, ok := raw["AcceptConnection"]; ok {
		var acceptConnection protocol.AcceptConnection
		if err := json.Unmarshal(data, &acceptConnection.UUID); err != nil {
			slog.Warn("incorrect AcceptConnection", "error", err)
			conn.Close()
			return
		}

		handleAcceptConnection(conn, acceptConnection)
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
