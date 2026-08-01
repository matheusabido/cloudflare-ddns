package config

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
var conflictiveTypes = map[CloudflareDDNSRecordType][]CloudflareDDNSRecordType{
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
		if current.Name == record.Name && slices.Contains(conflictiveTypes[current.Type], record.Type) {
			return fmt.Errorf("conflictive record types for name %s: %s and %s", record.Name, current.Type, record.Type)
		}
	}

	c.Records = append(c.Records, record)
	return nil
}
