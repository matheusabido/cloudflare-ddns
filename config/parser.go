package config

import (
	"fmt"
	"strings"

	"github.com/matheusabido/cloudflare-ddns/types"
	"github.com/matheusabido/cloudflare-ddns/utils"
)

// Parses the config and returns a CloudflareDDNSConfig instance.
// It returns a slice of errors if any errors were found during parsing.
func ParseConfig(lines []string) (*types.CloudflareDDNSConfig, []error) {
	config := types.NewConfig()

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
			config.APIToken = value
		case "zone_id":
			config.ZoneID = value
		case "last_ipv4":
			config.LastIPv4 = value
		case "last_ipv6":
			config.LastIPv6 = value
		case "record":
			record, err := ParseRecord(value)
			if err != nil {
				errors = append(errors, fmt.Errorf("invalid record line: %s. Error: %v", value, err))
				continue
			}

			if err := config.AddRecord(record); err != nil {
				errors = append(errors, err)
				continue
			}
		}
	}

	return config, errors
}

// Parses a record line and returns a CloudflareDDNSRecord instance.
// It returns an error if the line is incorrectly configured.
func ParseRecord(line string) (*types.CloudflareDDNSRecord, error) {
	parts := strings.Split(line, ",")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid record line: %s. Incomplete record config. Use: record=type,name,value,proxied?,ttl?", line)
	}

	recordType := types.CloudflareDDNSRecordType(parts[0])
	if !recordType.IsValid() {
		return nil, fmt.Errorf("invalid record type: %s. Supported types: %s", parts[0], types.GetSupportedRecordTypes())
	}

	proxied := true
	if len(parts) >= 4 {
		proxied = utils.IsActive(parts[3])
	}

	ttl := types.TTLAuto
	if len(parts) >= 5 {
		ttl = types.CloudflareDDNSTTL(parts[4])
		if !ttl.IsValid() {
			return nil, fmt.Errorf("invalid TTL: %s. Supported values: %s", parts[4], types.GetSupportedTTLs())
		}
	}

	return &types.CloudflareDDNSRecord{
		Type:    recordType,
		Name:    parts[1],
		Value:   parts[2],
		Proxied: proxied,
		TTL:     ttl,
	}, nil
}
