package config

import (
	"fmt"
	"strings"
)

func ParseConfig(lines []string) (*CloudflareDDNSConfig, []error) {
	config := NewConfig()

	errors := make([]error, 0)
	for _, line := range lines {
		index := strings.Index(line, "=")
		if index == -1 {
			errors = append(errors, fmt.Errorf("invalid line: %s. Missing equals sign", line))
			continue
		}

		key := strings.ToLower(strings.TrimSpace(line[:index]))
		value := strings.TrimSpace(line[index+1:])
		switch key {
		case "api_token":
			fmt.Printf("API TOKEN: %s\n", value)
		case "zone_id":
			fmt.Printf("ZONE ID: %s\n", value)
		case "record":
			fmt.Printf("RECORD: %s\n", value)
		}
	}

	return config, errors
}
