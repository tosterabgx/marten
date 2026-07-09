package server

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/tosterabgx/marten/internal/protocol"
)

var subdomainMu sync.RWMutex
var subdomainTunnels = make(map[string]net.Conn)

func RunHTTPServer(port uint16) error {
	addr := protocol.JoinAddr("localhost", port)
	slog.Info("listening http")
	return http.ListenAndServe(addr, http.HandlerFunc(handleHTTP))
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	subdomain := strings.SplitN(r.Host, ".", 2)[0]

	subdomainMu.RLock()
	controlConn, ok := subdomainTunnels[subdomain]
	subdomainMu.RUnlock()

	if !ok {
		slog.Warn("no tunnel for subdomain", "subdomain", subdomain)
		http.Error(w, "no active tunnel at this address", http.StatusNotFound)
		return
	}

	var reqBuf bytes.Buffer
	if err := r.Write(&reqBuf); err != nil {
		http.Error(w, "failed to serialize request", http.StatusBadGateway)
		return
	}

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		slog.Error("hijack not supported by this server")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	conn, _, err := hijacker.Hijack()
	if err != nil {
		slog.Error("hijack failed", "error", err)
		return
	}

	id := uuid.New()

	connsMu.Lock()
	externalConns[id] = pendingConn{conn: conn, initial: reqBuf.Bytes()}
	connsMu.Unlock()

	msg := protocol.NewMessage(protocol.NewConnection{UUID: id})
	if err := json.NewEncoder(controlConn).Encode(msg); err != nil {
		slog.Warn("failed to send NewConnection", "error", err, "subdomain", subdomain)
		conn.Close()

		connsMu.Lock()
		delete(externalConns, id)
		connsMu.Unlock()
		return
	}

	slog.Info("sent NewConnection", "subdomain", subdomain, "uuid", id)
}
