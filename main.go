package main

import (
	"fmt"

	"github.com/matheusabido/cloudflare-ddns/config"
)

func main() {
	lines := config.ReadConfigFile()
	ddns, errors := config.ParseConfig(lines)
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Printf("Error parsing config: %v\n", err)
		}
		return
	}

	fmt.Printf("Parsed config: %+v\n", ddns)
}
