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

	switch msg.Type {
	case protocol.TypeClientHello:
		var clientHello protocol.ClientHello
		if err := msg.Decode(&clientHello); err != nil {
			slog.Warn("incorrect ClientHello", "error", err)
			conn.Close()
			return
		}

		l, err := handleNewClient(conn, clientHello)
		if err != nil {
			slog.Warn("failed to handle ClientHello", "error", err)
			conn.Close()
			return
		}
		defer l.Close()

		io.Copy(io.Discard, conn)
	case protocol.TypeAcceptConnection:
		var acceptConnection protocol.AcceptConnection
		if err := msg.Decode(&acceptConnection); err != nil {
			slog.Warn("incorrect AcceptConnection", "error", err)
			conn.Close()
			return
		}

		err := handleAcceptConnection(conn, acceptConnection)
		if err != nil {
			slog.Warn("failed to handle AcceptConnection", "error", err)
			conn.Close()
		}
	default:
		slog.Warn("unknown message", "message", msg)
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
