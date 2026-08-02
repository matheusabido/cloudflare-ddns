package types

import (
	"fmt"
	"slices"
)

type CloudflareDDNSConfig struct {
	APIToken string
	ZoneID   string
	LastIPv4 string
	LastIPv6 string
	Records  []*CloudflareDDNSRecord
}

type CloudflareDDNSRecord struct {
	Type    CloudflareDDNSRecordType
	Name    string
	Value   string
	Proxied bool
	TTL     CloudflareDDNSTTL
}

func (c *CloudflareDDNSRecord) ParseName(zoneName string) any {
	if c.Name == "@" {
		return zoneName
	}
	return fmt.Sprintf("%s.%s", c.Name, zoneName)
}

type CloudflareDDNSTTL string

const (
	TTL1Min  CloudflareDDNSTTL = "1min"
	TTL2Min  CloudflareDDNSTTL = "2min"
	TTL5Min  CloudflareDDNSTTL = "5min"
	TTL10Min CloudflareDDNSTTL = "10min"
	TTL15Min CloudflareDDNSTTL = "15min"
	TTL30Min CloudflareDDNSTTL = "30min"
	TTL1H    CloudflareDDNSTTL = "1h"
	TTL2H    CloudflareDDNSTTL = "2h"
	TTL5H    CloudflareDDNSTTL = "5h"
	TTL12H   CloudflareDDNSTTL = "12h"
	TTL1D    CloudflareDDNSTTL = "1d"
	TTLAuto  CloudflareDDNSTTL = "auto"
)

func (t CloudflareDDNSTTL) GetValue() int {
	switch t {
	case TTL1Min:
		return 60
	case TTL2Min:
		return 120
	case TTL5Min:
		return 300
	case TTL10Min:
		return 600
	case TTL15Min:
		return 900
	case TTL30Min:
		return 1800
	case TTL1H:
		return 3600
	case TTL2H:
		return 7200
	case TTL5H:
		return 18000
	case TTL12H:
		return 43200
	case TTL1D:
		return 86400
	case TTLAuto:
		return 1
	default:
		return 1
	}
}

// Returns all valid TTL values
func GetSupportedTTLs() []CloudflareDDNSTTL {
	return []CloudflareDDNSTTL{
		TTL1Min,
		TTL2Min,
		TTL5Min,
		TTL10Min,
		TTL15Min,
		TTL30Min,
		TTL1H,
		TTL2H,
		TTL5H,
		TTL12H,
		TTL1D,
		TTLAuto,
	}
}

// Returns if the TTL is valid
func (t CloudflareDDNSTTL) IsValid() bool {
	return slices.Contains(GetSupportedTTLs(), t)
}

type CloudflareDDNSRecordType string

const (
	RecordTypeA     CloudflareDDNSRecordType = "A"
	RecordTypeAAAA  CloudflareDDNSRecordType = "AAAA"
	RecordTypeCNAME CloudflareDDNSRecordType = "CNAME"
	RecordTypeTXT   CloudflareDDNSRecordType = "TXT"
)

// Returns if the record type is valid
func (r CloudflareDDNSRecordType) IsValid() bool {
	return slices.Contains(GetSupportedRecordTypes(), r)
}

// Returns all supported record types
func GetSupportedRecordTypes() []CloudflareDDNSRecordType {
	return []CloudflareDDNSRecordType{
		RecordTypeA,
		RecordTypeAAAA,
		RecordTypeCNAME,
		RecordTypeTXT,
	}
}

// NewConfig creates a new CloudflareDDNSConfig instance with an empty list of records.
func NewConfig() *CloudflareDDNSConfig {
	return &CloudflareDDNSConfig{
		Records: make([]*CloudflareDDNSRecord, 0),
	}
}

// Key: supported record type
// Value: record types incompatible with the key record type
var ConflictiveTypes = map[CloudflareDDNSRecordType][]CloudflareDDNSRecordType{
	RecordTypeA:     {RecordTypeCNAME},
	RecordTypeAAAA:  {RecordTypeCNAME},
	RecordTypeCNAME: {RecordTypeA, RecordTypeAAAA},
	RecordTypeTXT:   {},
}

// Adds a record to the list. It checks if there's any conflict before doing so and returns an error if there is.
func (c *CloudflareDDNSConfig) AddRecord(record *CloudflareDDNSRecord) error {
	if record == nil {
		return nil
	}

	for _, current := range c.Records {
		if current.Name == record.Name && slices.Contains(ConflictiveTypes[current.Type], record.Type) {
			return fmt.Errorf("conflictive record types for name %s: %s and %s", record.Name, current.Type, record.Type)
		}
	}

	c.Records = append(c.Records, record)
	return nil
}
