package main

import (
	"fmt"
	"time"

	"github.com/matheusabido/cloudflare-ddns/config"
)

func main() {
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

	fmt.Println(ddns)
}
