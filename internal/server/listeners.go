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
	for i := desiredPort; i < 65535; i++ {
		portsMu.RLock()
		_, ok := occupiedPorts[i]
		portsMu.RUnlock()
		if !ok {
			return i, nil
		}
	}

	// TODO: return an error
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

func handleConnection(conn net.Conn, controlConn net.Conn) {
	defer conn.Close()

	uuid := uuid.New()

	outboundConnections[uuid] = conn

	newConnectinonRequest := protocol.NewConnection{UUID: uuid}
	enc := json.NewEncoder(controlConn)
	if err := enc.Encode(newConnectinonRequest); err != nil {
		fmt.Println("couldn't encode a NewConnections:", err)
		return
	}

	fmt.Printf("Sent NewConnection: %v\n", newConnectinonRequest)
}

func handleListener(l net.Listener, controlConn net.Conn) {
	defer l.Close()
	defer controlConn.Close()

	fmt.Println("start listening")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		go handleConnection(conn, controlConn)
	}
}
