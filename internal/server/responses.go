package server

type StatusResponse struct {
	Status        string `json:"status"`
	Uptime        string `json:"uptime"`
	ActiveTunnels int    `json:"active_tunnels"`
}
