package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/matheusabido/cloudflare-ddns/types"
)

// Updates the config file to reflect the new IP addresses. The rest is unchanged to preserve comments.
// It returns an error if it fails to read or write the file.
func UpdateConfigIP(config *types.CloudflareDDNSConfig) error {
	configurationPath := GetConfigPath()
	file, err := os.OpenFile(configurationPath, os.O_RDWR, 0640)
	if err != nil {
		return fmt.Errorf("Could not open file %s. %v", configurationPath, err)
	}
	defer file.Close()

	var builder strings.Builder
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return fmt.Errorf("Could not read file %s. %v", configurationPath, err)
		}

		if strings.HasPrefix(line, "last_ipv4=") {
			fmt.Fprintf(&builder, "last_ipv4=%s\n", config.LastIPv4)
			continue
		}

		if strings.HasPrefix(line, "last_ipv6=") {
			fmt.Fprintf(&builder, "last_ipv6=%s\n", config.LastIPv6)
			continue
		}

		builder.WriteString(line)

		if err == io.EOF {
			break
		}
	}

	err = os.WriteFile(configurationPath, []byte(builder.String()), 0640)
	if err != nil {
		return fmt.Errorf("Could not write to file %s. %v", configurationPath, err)
	}
	return nil
}
