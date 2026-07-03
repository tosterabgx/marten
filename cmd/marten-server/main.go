package main

import (
	"log"
	"log/slog"

	"github.com/tosterabgx/marten/internal/protocol"
	"github.com/tosterabgx/marten/internal/server"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	slog.Info("starting Marten server")

	go func() {
		if err := server.RunHTTPServer(protocol.HTTPPort); err != nil {
			log.Fatal("http server failed:", err)
		}
	}()

	if err := server.RunControlServer(protocol.ControlPort); err != nil {
		log.Fatal("control server failed:", err)
	}
}
