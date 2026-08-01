package utils

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type NetworkType string

const (
	NetworkTypeIPv4 NetworkType = "tcp4"
	NetworkTypeIPv6 NetworkType = "tcp6"
)

// GetIP fetches the public IP address of the machine using the specified network protocol, either "tcp4" or "tcp6".
// It returns the IP address as a string or an error encountered during the process.
func GetIP(network NetworkType) (string, error) {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, _, addr string) (net.Conn, error) {
			d := &net.Dialer{Timeout: 5 * time.Second}
			return d.DialContext(ctx, string(network), addr) // "tcp4" ou "tcp6"
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	req, err := http.NewRequest("GET", "https://ifconfig.me", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "curl/7.88.1")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ip := strings.TrimSpace(string(body))

	if ip == "" {
		return "", fmt.Errorf("could not fetch IP address, response body is empty")
	}

	return ip, nil
}
