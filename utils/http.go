package utils

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/matheusabido/cloudflare-ddns/types"
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
			return d.DialContext(ctx, string(network), addr)
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

// FetchIP fetches the current active IP addresses
// It also checks if they have changed compared to the last known values and updates the config (in memory) accordingly
func FetchIP(ddns *types.CloudflareDDNSConfig) (bool, error) {
	ipChanged := false
	if IsActive(ddns.LastIPv4) {
		ipv4, err := GetIP(NetworkTypeIPv4)
		if err != nil {
			return false, fmt.Errorf("Error getting IPv4: %v\n", err)
		}

		if ipv4 != ddns.LastIPv4 {
			fmt.Printf("IPv4 has changed from \"%s\" to \"%s\"\n", ddns.LastIPv4, ipv4)
			ddns.LastIPv4 = ipv4
			ipChanged = true
		}
	}

	if IsActive(ddns.LastIPv6) {
		ipv6, err := GetIP(NetworkTypeIPv6)
		if err != nil {
			return false, fmt.Errorf("Error getting IPv6: %v\n", err)
		}

		if ipv6 != ddns.LastIPv6 {
			fmt.Printf("IPv6 has changed from \"%s\" to \"%s\"\n", ddns.LastIPv6, ipv6)
			ddns.LastIPv6 = ipv6
			ipChanged = true
		}
	}
	return ipChanged, nil
}
