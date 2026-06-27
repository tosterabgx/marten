package main

import (
	"log/slog"

	"github.com/tosterabgx/marten/internal/server"
)

func main() {
	slog.Info("starting Marten server")
	if err := server.RunControlServer(); err != nil {
		panic(err)
	}
}
