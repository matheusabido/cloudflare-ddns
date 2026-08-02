package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/matheusabido/cloudflare-ddns/config"
	"github.com/matheusabido/cloudflare-ddns/tools"
	"github.com/matheusabido/cloudflare-ddns/utils"
)

func main() {
	forceUpdate := flag.Bool("force", false, "Forces update regardless of the ip having changed")
	init := flag.Bool("init", false, "Creates the default configuration file")
	flag.Parse()

	if *init {
		fmt.Println("Creating the default config file.")
		if err := config.WriteDefaultConfigFile(); err != nil {
			log.Panicf("Error creating default config file: %v", err)
		}
		fmt.Println("Default config created.")
		return
	}

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
	ipChanged, err := utils.FetchIP(ddns)
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
