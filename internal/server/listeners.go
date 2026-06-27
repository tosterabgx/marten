package server

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/google/uuid"
	"github.com/tosterabgx/marten/internal/protocol"
)

var portsMu sync.RWMutex
var occupiedPorts = make(map[uint16]struct{})
var connsMu sync.RWMutex
var outboundConnections = make(map[uuid.UUID]net.Conn)

func getAvailablePort(desiredPort uint16) (uint16, error) {
	for i := protocol.MinPort; i < protocol.MaxPort; i++ {
		portsMu.RLock()
		_, ok := occupiedPorts[i]
		portsMu.RUnlock()
		if !ok {
			return i, nil
		}
	}

	return 0, nil
}

func registerListener(desiredPort uint16, controlConn net.Conn) (uint16, error) {
	port, _ := getAvailablePort(desiredPort)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return 0, err
	}

	portsMu.Lock()
	occupiedPorts[port] = struct{}{}
	portsMu.Unlock()

	go handleListener(l, controlConn)

	return port, nil
}

func handleTunnelConnection(conn net.Conn, controlConn net.Conn) {
	uuid := uuid.New()

	connsMu.Lock()
	outboundConnections[uuid] = conn
	connsMu.Unlock()

	newConnectinonRequest := protocol.NewConnection{UUID: uuid}
	enc := json.NewEncoder(controlConn)
	if err := enc.Encode(newConnectinonRequest); err != nil {
		fmt.Println("Error encoding NewConnection:", err)
		conn.Close()
		return
	}

	fmt.Println("Sent NewConnection:", newConnectinonRequest)
}

func handleListener(l net.Listener, controlConn net.Conn) {
	defer l.Close()
	defer controlConn.Close()

	fmt.Println("Start listening")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		fmt.Println("Got a connection")

		go handleTunnelConnection(conn, controlConn)
	}
}
