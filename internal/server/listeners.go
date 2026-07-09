package server

import (
	"errors"
	"log/slog"
	"math/rand"
	"net"
	"sync"

	"github.com/tosterabgx/marten/internal/protocol"
)

var portsMu sync.RWMutex
var occupiedPorts = make(map[uint16]struct{})

func getAvailablePort() (uint16, error) {
	const maxAttempts = 150

	rangeSize := int(protocol.MaxPort) - int(protocol.MinPort) + 1

	for attempt := 0; attempt < maxAttempts; attempt++ {
		port := protocol.MinPort + uint16(rand.Intn(rangeSize))

		portsMu.Lock()
		_, ok := occupiedPorts[port]
		if !ok {
			occupiedPorts[port] = struct{}{}
		}
		portsMu.Unlock()

		if !ok {
			return port, nil
		}
	}

	return 0, errors.New("no available port after 150 attempts")
}

func registerListener(controlConn net.Conn) (net.Listener, uint16, error) {
	port, err := getAvailablePort()
	if err != nil {
		return nil, 0, err
	}

	addr := protocol.JoinAddr("", port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		portsMu.Lock()
		delete(occupiedPorts, port)
		portsMu.Unlock()
		return nil, 0, err
	}

	go handleListener(l, controlConn, port)

	return l, port, nil
}

func handleListener(l net.Listener, controlConn net.Conn, port uint16) {
	defer cleanupListener(l, port)

	slog.Info("listening external", "address", l.Addr().String())

	for {
		conn, err := l.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}

			slog.Warn("accept failed", "error", err)
			continue
		}

		go handleExternalConnection(conn, controlConn)
	}
}

func cleanupListener(l net.Listener, port uint16) {
	l.Close()

	portsMu.Lock()
	delete(occupiedPorts, port)
	portsMu.Unlock()

	slog.Debug("listener cleaned up", "port", port)
}
