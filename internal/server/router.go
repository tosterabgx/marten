package server

import (
	"encoding/json"
	"log/slog"
	"net"
	"sync"

	"github.com/google/uuid"
	"github.com/tosterabgx/marten/internal/protocol"
)

type pendingConn struct {
	conn    net.Conn
	initial []byte
}

var connsMu sync.RWMutex
var externalConns = make(map[uuid.UUID]pendingConn)

func handleExternalConnection(conn net.Conn, controlConn net.Conn) {
	slog.Debug("got new external connection")

	uuid := uuid.New()

	connsMu.Lock()
	externalConns[uuid] = pendingConn{conn: conn}
	connsMu.Unlock()

	msg := protocol.NewMessage(protocol.NewConnection{UUID: uuid})

	if err := json.NewEncoder(controlConn).Encode(msg); err != nil {
		slog.Warn("NewConnection encoding failed", "error", err)

		connsMu.Lock()
		delete(externalConns, uuid)
		connsMu.Unlock()

		conn.Close()
		return
	}

	slog.Debug("sent NewConnection", "message", msg)
}
