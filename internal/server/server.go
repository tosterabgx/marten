package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net"

	"github.com/tosterabgx/marten/internal/protocol"
)

func handleConnection(conn net.Conn) {
	dec := json.NewDecoder(conn)
	dec.DisallowUnknownFields()

	var msg protocol.Message
	if err := dec.Decode(&msg); err != nil {
		slog.Warn("incorrect message", "error", err)
		conn.Close()
		return
	}

	switch p := msg.Payload.(type) {
	case protocol.ClientHello:
		if p.Type == "http" {
			subdomain, err := GetAvailableName(&conn)
			if err != nil {
				conn.Close()
				return
			}

			serverHello := protocol.NewMessage(protocol.ServerHello{Subdomain: subdomain})
			if err := json.NewEncoder(conn).Encode(serverHello); err != nil {
				slog.Warn("ServerHello encoding failed", "error", err)
			}
			return
		}
		l, err := handleNewClient(conn, p)
		if err != nil {
			slog.Warn("failed to handle ClientHello", "error", err)
			conn.Close()
			return
		}
		defer l.Close()
		io.Copy(io.Discard, conn)

	case protocol.AcceptConnection:
		err := handleAcceptConnection(conn, p)
		if err != nil {
			slog.Warn("failed to handle AcceptConnection", "error", err)
			conn.Close()
		}

	default:
		slog.Warn("unknown message", "message", msg)
	}
}

func RunControlServer(port uint16) error {
	addr := protocol.JoinAddr("", port)
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
