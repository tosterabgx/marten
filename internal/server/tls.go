package server

import (
	"crypto/tls"
	"fmt"
	"os"
)

func LoadTLSConfig() (*tls.Config, error) {
	certFile := os.Getenv("MARTEN_TLS_CERT")
	keyFile := os.Getenv("MARTEN_TLS_KEY")

	if certFile == "" {
		certFile = "/var/lib/caddy/.local/share/caddy/certificates/acme-v02.api.letsencrypt.org-directory/wildcard_.usemarten.tech/wildcard_.usemarten.tech.crt"
	}
	if keyFile == "" {
		keyFile = "/var/lib/caddy/.local/share/caddy/certificates/acme-v02.api.letsencrypt.org-directory/wildcard_.usemarten.tech/wildcard_.usemarten.tech.key"
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("loading TLS cert: %v", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
	}, nil
}
