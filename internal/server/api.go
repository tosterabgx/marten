package server

import (
	"net/http"

	"github.com/tosterabgx/marten/internal/protocol"
)

func RunAPIServer(port uint16) error {
	addr := protocol.JoinAddr("localhost", port)

	http.HandleFunc("/api/status", statusHandler)
	return http.ListenAndServe(addr, nil)
}
