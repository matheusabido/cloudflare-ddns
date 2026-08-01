package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"

	"github.com/matheusabido/cloudflare-ddns/types"
)

type CloudflareClient struct {
	config  *types.CloudflareDDNSConfig
	client  *http.Client
	baseURL string
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
	recordsList, err := c.ListRecords()
	if err != nil {
		return err
	}

	for _, record := range c.config.Records {
		conflictiveTypes := types.ConflictiveTypes[record.Type]
		equivalentId := ""
		hasConflict := false

		for _, cfRecord := range recordsList.Result {
			isEquivalent := record.Name == cfRecord.Name && record.Type == types.CloudflareDDNSRecordType(cfRecord.Type)
			if isEquivalent {
				equivalentId = cfRecord.ID
				continue
			}

			if slices.Contains(conflictiveTypes, types.CloudflareDDNSRecordType(cfRecord.Type)) {
				fmt.Printf("Found conflictive record on Cloudflare: [%s %s] with ID: %s\n", cfRecord.Type, cfRecord.Name, cfRecord.ID)
				hasConflict = true
				break
			}
		}
		if hasConflict {
			fmt.Printf("DDNS Record has conflicts with existing Cloudflare records. Skipping it. [%s %s]", record.Type, record.Name)
			continue
		}

		if equivalentId != "" {
			if err := c.OverwriteRecord(equivalentId, record); err != nil {
				return err
			}
		} else {
			if err := c.CreateRecord(record); err != nil {
				return err
			}
		}
	}
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

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error while trying to read record list response: %v", err)
	}

	var response types.CloudflareListRecordsResponse
	if err := json.Unmarshal(content, &response); err != nil {
		return nil, fmt.Errorf("Error while trying to parse record list response: %v", err)
	}

	return &response, nil
}

/*
	{
	  "errors": [
	    {
	      "code": 1000,
	      "message": "message",
	      "documentation_url": "documentation_url",
	      "source": {
	        "pointer": "pointer"
	      }
	    }
	  ],
	  "messages": [
	    {
	      "code": 1000,
	      "message": "message",
	      "documentation_url": "documentation_url",
	      "source": {
	        "pointer": "pointer"
	      }
	    }
	  ],
	  "success": true,
	  "result": {
	    "name": "example.com",
	    "ttl": 3600,
	    "type": "A",
	    "comment": "Domain verification record",
	    "content": "198.51.100.4",
	    "private_routing": true,
	    "proxied": true,
	    "settings": {
	      "ipv4_only": true,
	      "ipv6_only": true
	    },
	    "tags": [
	      "owner:dns-team"
	    ],
	    "id": "023e105f4ecef8ad9ca31a8372d0c353",
	    "created_on": "2014-01-01T05:20:00.12345Z",
	    "meta": {
	      "dead_glue": true,
	      "is_glue": true,
	      "shadowed_by": [
	        "372e67954025e0ba6aaa6d586b9e0b59"
	      ],
	      "shadowed_records_count": 42
	    },
	    "modified_on": "2014-01-01T05:20:00.12345Z",
	    "proxiable": true,
	    "comment_modified_on": "2024-01-01T05:20:00.12345Z",
	    "tags_modified_on": "2025-01-01T05:20:00.12345Z"
	  }
	}

	curl https://api.cloudflare.com/client/v4/zones/$ZONE_ID/dns_records \
	    -H 'Content-Type: application/json' \
	    -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
	    -d '{
	          "name": "example.com",
	          "ttl": 3600,
	          "type": "A",
	          "comment": "Domain verification record",
	          "content": "198.51.100.4",
	          "private_routing": true,
	          "proxied": true
	        }'
*/
func (c *CloudflareClient) CreateRecord(record *types.CloudflareDDNSRecord) error {
	return nil
}

/*
	curl https://api.cloudflare.com/client/v4/zones/$ZONE_ID/dns_records/$DNS_RECORD_ID \
	    -X PUT \
	    -H 'Content-Type: application/json' \
	    -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
	    -d '{
	          "name": "example.com",
	          "ttl": 3600,
	          "type": "A",
	          "comment": "Domain verification record",
	          "content": "198.51.100.4",
	          "private_routing": true,
	          "proxied": true
	        }'

			{
	  "errors": [
	    {
	      "code": 1000,
	      "message": "message",
	      "documentation_url": "documentation_url",
	      "source": {
	        "pointer": "pointer"
	      }
	    }
	  ],
	  "messages": [
	    {
	      "code": 1000,
	      "message": "message",
	      "documentation_url": "documentation_url",
	      "source": {
	        "pointer": "pointer"
	      }
	    }
	  ],
	  "success": true,
	  "result": {
	    "name": "example.com",
	    "ttl": 3600,
	    "type": "A",
	    "comment": "Domain verification record",
	    "content": "198.51.100.4",
	    "private_routing": true,
	    "proxied": true,
	    "settings": {
	      "ipv4_only": true,
	      "ipv6_only": true
	    },
	    "tags": [
	      "owner:dns-team"
	    ],
	    "id": "023e105f4ecef8ad9ca31a8372d0c353",
	    "created_on": "2014-01-01T05:20:00.12345Z",
	    "meta": {
	      "dead_glue": true,
	      "is_glue": true,
	      "shadowed_by": [
	        "372e67954025e0ba6aaa6d586b9e0b59"
	      ],
	      "shadowed_records_count": 42
	    },
	    "modified_on": "2014-01-01T05:20:00.12345Z",
	    "proxiable": true,
	    "comment_modified_on": "2024-01-01T05:20:00.12345Z",
	    "tags_modified_on": "2025-01-01T05:20:00.12345Z"
	  }
	}
*/
func (c *CloudflareClient) OverwriteRecord(id string, record *types.CloudflareDDNSRecord) error {
	return nil
}
