package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/matheusabido/cloudflare-ddns/types"
	"github.com/matheusabido/cloudflare-ddns/utils"
)

type CloudflareClient struct {
	config   *types.CloudflareDDNSConfig
	client   *http.Client
	baseURL  string
	zoneName string
}

func NewCloudflareClient(config *types.CloudflareDDNSConfig) *CloudflareClient {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &CloudflareClient{
		config:  config,
		client:  client,
		baseURL: fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s", config.ZoneID),
	}
}

// UpdateRecords updates the DNS records in Cloudflare based on the current configuration.
// Returns an error if the update fails.
func (c *CloudflareClient) UpdateRecords() error {
	start := time.Now()

	fmt.Println("Fetching zone details...")
	startDetails := time.Now()
	zoneDetails, err := c.GetZoneDetails()
	if err != nil {
		return err
	}

	if !zoneDetails.Success {
		zoneBytes, err := json.Marshal(zoneDetails)
		if err != nil {
			return fmt.Errorf("could not get cloudflare zone details. Could not marshal response: %v", err)
		}
		return fmt.Errorf("could not get cloudflare zone details. Response: %v", string(zoneBytes))
	}
	c.zoneName = zoneDetails.Result.Name
	fmt.Printf("Fetched cloudflare zone details in %s\n", time.Since(startDetails))

	startRecords := time.Now()
	recordsList, err := c.ListRecords()
	if err != nil {
		return err
	}
	fmt.Printf("Listed cloudflare records in %s\n", time.Since(startRecords))

	startUpdate := time.Now()
	fmt.Println("Updating cloudflare...")
	ipv4Active := utils.IsActive(c.config.LastIPv4)
	ipv6Active := utils.IsActive(c.config.LastIPv6)
	for _, record := range c.config.Records {
		if !ipv4Active && record.Type == types.RecordTypeA {
			fmt.Printf("Skipping %s %s because IPv4 is not active\n", record.Type, record.Name)
			continue
		}

		if !ipv6Active && record.Type == types.RecordTypeAAAA {
			fmt.Printf("Skipping %s %s because IPv6 is not active\n", record.Type, record.Name)
			continue
		}
		startConflict := time.Now()
		fmt.Printf("Checking for conflicts for %s %s...\n", record.Type, record.Name)

		conflictiveTypes := types.ConflictiveTypes[record.Type]
		var equivalentRecord *types.CloudflareRecord = nil
		hasConflict := false

		recordName := record.ParseName(zoneDetails.Result.Name)
		for _, cfRecord := range recordsList.Result {
			isEquivalent := recordName == cfRecord.Name && record.Type == types.CloudflareDDNSRecordType(cfRecord.Type)
			if isEquivalent {
				equivalentRecord = &cfRecord
				continue
			}

			if slices.Contains(conflictiveTypes, types.CloudflareDDNSRecordType(cfRecord.Type)) {
				fmt.Printf("Found conflictive record on Cloudflare: [%s %s] with ID: %s\n", cfRecord.Type, cfRecord.Name, cfRecord.ID)
				hasConflict = true
				break
			}
		}
		fmt.Printf("Checked for conflicts in %s\n", time.Since(startConflict))

		if hasConflict {
			fmt.Printf("DDNS Record has conflicts with existing Cloudflare records. Skipping it. [%s %s]", record.Type, record.Name)
			continue
		}

		startRecord := time.Now()
		if equivalentRecord != nil {
			if err := c.OverwriteRecord(equivalentRecord, record); err != nil {
				return err
			}
			fmt.Printf("Overwrote %s %s in %s\n", record.Type, record.Name, time.Since(startRecord))
		} else {
			if err := c.CreateRecord(record); err != nil {
				return err
			}
			fmt.Printf("Created %s %s in %s\n", record.Type, record.Name, time.Since(startRecord))
		}
	}
	fmt.Printf("Cloudflare update done in %s\n", time.Since(startUpdate))
	fmt.Printf("All Cloudflare interactions done in %s\n", time.Since(start))
	return nil
}

// ListRecords fetches the list of DNS records from Cloudflare for the configured zone.
// Returns a CloudflareListRecordsResponse containing the records and any errors encountered.
func (c *CloudflareClient) ListRecords() (*types.CloudflareListRecordsResponse, error) {
	request, err := http.NewRequest("GET", c.baseURL+"/dns_records?per_page=5000000", nil)
	if err != nil {
		return nil, fmt.Errorf("Error while trying to create request to list cloudflare records: %v", err)
	}

	request.Header.Set("Authorization", "Bearer "+c.config.APIToken)
	resp, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Error while trying to fetch cloudflare records: %v", err)
	}
	defer resp.Body.Close()

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error while trying to read record list response: %v", err)
	}

	var response types.CloudflareListRecordsResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return nil, fmt.Errorf("Error while trying to parse record list response: %v", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("could not list cloudflare records. Response: %s", string(responseBytes))
	}

	return &response, nil
}

// Fetches the details of the configured Cloudflare zone.
// Returns a CloudflareZoneDetails containing the zone information and any errors encountered.
func (c *CloudflareClient) GetZoneDetails() (*types.CloudflareZoneDetails, error) {
	request, err := http.NewRequest("GET", c.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error while trying to create request to get cloudflare zone details: %v", err)
	}

	request.Header.Set("Authorization", "Bearer "+c.config.APIToken)

	resp, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Error while trying to fetch cloudflare zone details: %v", err)
	}
	defer resp.Body.Close()

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error while trying to read zone details response: %v", err)
	}

	var response types.CloudflareZoneDetails
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return nil, fmt.Errorf("Error while trying to parse zone details response: %v", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("could not get cloudflare zone details. Response: %s", string(responseBytes))
	}
	return &response, nil
}

// Creates a record on Cloudflare based on the provided CloudflareDDNSRecord.
// Returns an error if the creation fails.
func (c *CloudflareClient) CreateRecord(record *types.CloudflareDDNSRecord) error {
	recordValue := strings.TrimSpace(record.Value)
	recordValue = strings.ReplaceAll(recordValue, "{public_ipv4}", c.config.LastIPv4)
	recordValue = strings.ReplaceAll(recordValue, "{public_ipv6}", c.config.LastIPv6)
	body := map[string]any{
		"name":    strings.TrimSpace(record.Name),
		"ttl":     record.TTL.GetValue(),
		"type":    string(record.Type),
		"content": recordValue,
		"proxied": record.Proxied,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("could not marshal create record body. %v", err)
	}

	buffer := bytes.NewBuffer(bodyBytes)
	request, err := http.NewRequest("POST", c.baseURL+"/dns_records", buffer)
	if err != nil {
		return fmt.Errorf("error while trying to create request to list cloudflare records: %v", err)
	}

	request.Header.Set("Authorization", "Bearer "+c.config.APIToken)

	resp, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf("could not create record on Cloudflare: %v", err)
	}
	defer resp.Body.Close()

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read create record response body: %v", err)
	}

	var response types.CloudflareCreateRecordResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return fmt.Errorf("error while trying to parse create record response: %v", err)
	}

	if !response.Success {
		return fmt.Errorf("could not create cloudflare record. Response: %s", string(responseBytes))
	}
	return nil
}

// Updates an existing DNS record on Cloudflare with the provided ID and CloudflareDDNSRecord.
// Returns an error if the update fails.
func (c *CloudflareClient) OverwriteRecord(equivalentRecord *types.CloudflareRecord, record *types.CloudflareDDNSRecord) error {
	recordValue := strings.TrimSpace(record.Value)
	recordValue = strings.ReplaceAll(recordValue, "{public_ipv4}", c.config.LastIPv4)
	recordValue = strings.ReplaceAll(recordValue, "{public_ipv6}", c.config.LastIPv6)

	if equivalentRecord.Content == record.ParseName(c.zoneName) && equivalentRecord.Proxied == record.Proxied && equivalentRecord.TTL == record.TTL.GetValue() {
		fmt.Printf("No changes detected for %s %s. Skipping update.\n", record.Type, record.Name)
		return nil
	}

	body := map[string]any{
		"name":    strings.TrimSpace(record.Name),
		"ttl":     record.TTL.GetValue(),
		"type":    string(record.Type),
		"content": recordValue,
		"proxied": record.Proxied,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("could not marshal overwrite record body. %v", err)
	}

	buffer := bytes.NewBuffer(bodyBytes)
	request, err := http.NewRequest("PUT", c.baseURL+"/dns_records/"+equivalentRecord.ID, buffer)

	request.Header.Set("Authorization", "Bearer "+c.config.APIToken)

	resp, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf("could not overwrite record on Cloudflare: %v", err)
	}
	defer resp.Body.Close()

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read overwrite record response body: %v", err)
	}

	var response types.CloudflareCreateRecordResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return fmt.Errorf("error while trying to parse overwrite record response: %v", err)
	}

	if !response.Success {
		return fmt.Errorf("could not overwrite cloudflare record. Response: %s", string(responseBytes))
	}
	return nil
}
