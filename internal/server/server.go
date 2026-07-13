package server

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"time"

	"github.com/tosterabgx/marten/internal/protocol"
)

var startTime = time.Now()

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
		if p.Type == protocol.TypeHTTP {
			subdomain, err := GetAvailableName(&conn)
			if err != nil {
				slog.Warn("failed to handle ClientHello", "error", err)
				conn.Close()
				return
			}

			serverHello := protocol.NewMessage(protocol.ServerHello{Subdomain: subdomain})
			if err := json.NewEncoder(conn).Encode(serverHello); err != nil {
				slog.Warn("ServerHello encoding failed", "error", err)
			}

			defer func() {
				subdomainMu.Lock()
				delete(subdomainTunnels, subdomain)
				subdomainMu.Unlock()
			}()
			io.Copy(io.Discard, conn)

		} else if p.Type == protocol.TypeTCP {
			l, err := handleNewClient(conn, p)
			if err != nil {
				slog.Warn("failed to handle ClientHello", "error", err)
				conn.Close()
				return
			}
			defer l.Close()
			io.Copy(io.Discard, conn)

		} else {
			slog.Warn("unknown ClientHello type", "type", p.Type)
			conn.Close()
			return
		}

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

	tlsConfig, err := LoadTLSConfig()
	if err != nil {
		return fmt.Errorf("TLS setup: %w", err)
	}

	l, err := tls.Listen("tcp", addr, tlsConfig)
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
