package config

import (
	"bufio"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
)

// ReadConfigFile reads the configuration file on the path determined by GetConfigPath()
// it returns a slice of strings containing the lines of the configuration file,
// ignoring comments and empty lines. It panics on failure to open or read the file.
func ReadConfigFile() []string {
	configurationPath := GetConfigPath()
	file, err := os.OpenFile(configurationPath, os.O_RDWR, 0640)
	if err != nil {
		log.Panicf("Could not open %s. Check if the file exists and this program has the correct permissions.", configurationPath)
	}
	defer file.Close()

	file.Seek(0, io.SeekStart)
	reader := bufio.NewReader(file)

	lines := make([]string, 0)

	bytes, err := reader.ReadBytes('\n')
	for err == nil || err == io.EOF {
		value := strings.TrimSpace(string(bytes))
		if strings.HasPrefix(value, "#") || value == "" {
			if err == io.EOF {
				break
			}
			bytes, err = reader.ReadBytes('\n')
			continue
		}

		index := strings.Index(value, "#")
		if index != -1 {
			value = value[:index]
		}

		lines = append(lines, value)

		if err == io.EOF {
			break
		}
		bytes, err = reader.ReadBytes('\n')
	}

	if err != io.EOF {
		log.Panicf("An error occurred while reading the %s file", configurationPath)
	}
	return lines
}

// GetConfigPath returns the path to the configuration file based on the operating system.
func GetConfigPath() string {
	if runtime.GOOS == "windows" {
		return ".\\cloudflare-ddns.conf"
	} else {
		return "/etc/cloudflare-ddns.conf"
	}
}
