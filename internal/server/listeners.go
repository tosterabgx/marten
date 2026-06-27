package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"

	"github.com/google/uuid"
	"github.com/tosterabgx/marten/internal/protocol"
)

var portsMu sync.RWMutex
var occupiedPorts = make(map[uint16]struct{})
var connsMu sync.RWMutex
var externalConns = make(map[uuid.UUID]net.Conn)

func getAvailablePort() (uint16, error) {
	for i := protocol.MinPort; i < protocol.MaxPort; i++ {
		portsMu.RLock()
		_, ok := occupiedPorts[i]
		portsMu.RUnlock()
		if !ok {
			return i, nil
		}
	}

	return 0, errors.New("no available port")
}

func registerListener(controlConn net.Conn) (uint16, error) {
	port, err := getAvailablePort()
	if err != nil {
		return 0, err
	}

	addr := protocol.JoinAddr("", port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return 0, err
	}

	portsMu.Lock()
	occupiedPorts[port] = struct{}{}
	portsMu.Unlock()

	go handleListener(l, controlConn)

	return port, nil
}

func handleListener(l net.Listener, controlConn net.Conn) {
	defer l.Close()
	defer controlConn.Close()

	slog.Info("listening external", "address", l.Addr().String())

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		go handleExternalConnection(conn, controlConn)
	}
}

func handleExternalConnection(conn net.Conn, controlConn net.Conn) {
	slog.Debug("got new external connection")

	uuid := uuid.New()

	connsMu.Lock()
	externalConns[uuid] = conn
	connsMu.Unlock()

	newConnectionRequest := protocol.NewConnection{UUID: uuid}
	enc := json.NewEncoder(controlConn)
	if err := enc.Encode(newConnectionRequest); err != nil {
		slog.Warn("NewConnection encoding failed", "error", err)
		conn.Close()
		return
	}

	slog.Debug("sent NewConnection", "message", newConnectionRequest)
}
