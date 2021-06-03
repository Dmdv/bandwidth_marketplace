package http

import (
	"net"
	"net/http"
	"time"
)

const (
	// clientTimeout represents default http.Client timeout.
	clientTimeout = 10 * time.Second

	// tlsHandshakeTimeout represents default http.Transport TLS handshake timeout.
	tlsHandshakeTimeout = 5 * time.Second

	// dialTimeout represents default net.Dialer timeout.
	dialTimeout = 5 * time.Second
)

// NewClient creates default http.Client with timeouts.
func NewClient() *http.Client {
	return &http.Client{
		Timeout: clientTimeout,
		Transport: &http.Transport{
			TLSHandshakeTimeout: tlsHandshakeTimeout,
			DialContext: (&net.Dialer{
				Timeout: dialTimeout,
			}).DialContext,
		},
	}
}
