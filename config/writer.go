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

// Writes the default configuration file to the config path. It returns an error if it fails to write the file.
func WriteDefaultConfigFile() error {
	defaultConfigContent := `api_token=[APITOKEN]
zone_id=[ZONEID]

# an empty value for last_ipv4 and last_ipv6 will let the script update the records on the first run
# if you set the value to "none", "disabled" or "no", it's gonna be disabled
# if ipv4 is disabled, A records will be ignored
# if ipv6 is disabled, AAAA records will be ignored
last_ipv4=
last_ipv6=none

# supported types: A, AAAA, CNAME, TXT
# ttl supported values = 1min, 2min, 5min, 10min, 15min, 30min, 1h, 2h, 5h, 12h, 1d, auto
# accepted value vars: {public_ipv4} {public_ipv6}

# A record is defined like so:
# record=type,name,value,proxied?,ttl?
#
# proxied and ttl are optional
# proxied's default is true
# ttl's default is auto

# This is an example A record, with the name "ddns", pointing to your public_ipv4. Proxied is defaulted to true. TTL is defaulted to "auto".
record=A,ddns,{public_ipv4}

# In this example below, proxy is disabled
#record=A,ddns,{public_ipv4},false`

	configurationPath := GetConfigPath()
	err := os.WriteFile(configurationPath, []byte(defaultConfigContent), 0640)
	if err != nil {
		return fmt.Errorf("Could not write to file %s. %v", configurationPath, err)
	}
	return nil
}
