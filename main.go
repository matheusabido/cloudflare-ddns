package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/matheusabido/cloudflare-ddns/config"
	"github.com/matheusabido/cloudflare-ddns/tools"
	"github.com/matheusabido/cloudflare-ddns/types"
	"github.com/matheusabido/cloudflare-ddns/utils"
)

func main() {
	forceUpdate := flag.Bool("force", false, "Forces update regardless of the ip having changed")

	start := time.Now()
	lines := config.ReadConfigFile()
	fmt.Printf("Read config in %s\n", time.Since(start))

	startParsing := time.Now()
	ddns, errors := config.ParseConfig(lines)
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Printf("Error parsing config: %v\n", err)
		}
		return
	}
	fmt.Printf("Parsed config in %s. Total: %s\n", time.Since(startParsing), time.Since(start))

	startFetching := time.Now()
	ipChanged, err := FetchIP(ddns)
	if err != nil {
		log.Panicf("could not fetch IP. %v", err)
	}
	fmt.Printf("Fetched IPs in %s. Total: %s\n", time.Since(startFetching), time.Since(start))

	if ipChanged || *forceUpdate {
		cloudflare := tools.NewCloudflareClient(ddns)
		if err := cloudflare.UpdateRecords(); err != nil {
			fmt.Printf("Error updating records: %v\n", err)
		}

		if err := config.UpdateConfigIP(ddns); err != nil {
			fmt.Printf("Error updating config: %v\n", err)
		}
	}

	fmt.Println("Done in", time.Since(start))
}

// FetchIP fetches the current active IP addresses
// It also checks if they have changed compared to the last known values and updates the config (in memory) accordingly
func FetchIP(ddns *types.CloudflareDDNSConfig) (bool, error) {
	ipChanged := false
	if utils.IsActive(ddns.LastIPv4) {
		ipv4, err := utils.GetIP(utils.NetworkTypeIPv4)
		if err != nil {
			return false, fmt.Errorf("Error getting IPv4: %v\n", err)
		}

		if ipv4 != ddns.LastIPv4 {
			fmt.Printf("IPv4 has changed from \"%s\" to \"%s\"\n", ddns.LastIPv4, ipv4)
			ddns.LastIPv4 = ipv4
			ipChanged = true
		}
	}

	if utils.IsActive(ddns.LastIPv6) {
		ipv6, err := utils.GetIP(utils.NetworkTypeIPv6)
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
