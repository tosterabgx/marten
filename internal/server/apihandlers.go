package server

import (
	"encoding/json"
	"net/http"
	"time"
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
	subdomainMu.RLock()
	httpTunnels := len(subdomainTunnels)
	subdomainMu.RUnlock()

	portsMu.RLock()
	tcpTunnels := len(occupiedPorts)
	portsMu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(StatusResponse{
		Status:        "ok",
		Uptime:        time.Since(startTime).Round(time.Second).String(),
		ActiveTunnels: httpTunnels + tcpTunnels,
	})
}
