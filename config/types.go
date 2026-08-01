package config

import (
	"fmt"
	"slices"
)

type CloudflareDDNSConfig struct {
	APIToken string
	ZoneID   string
	Records  []CloudflareDDNSRecord
}

type CloudflareDDNSRecord struct {
	Type    CloudflareDDNSRecordType
	Name    string
	Value   string
	TTL     CloudflareDDNSTTL
	Proxied bool
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

type CloudflareDDNSRecordType string

const (
	RecordTypeA     CloudflareDDNSRecordType = "A"
	RecordTypeAAAA  CloudflareDDNSRecordType = "AAAA"
	RecordTypeCNAME CloudflareDDNSRecordType = "CNAME"
	RecordTypeTXT   CloudflareDDNSRecordType = "TXT"
)

func NewConfig() *CloudflareDDNSConfig {
	return &CloudflareDDNSConfig{
		Records: make([]CloudflareDDNSRecord, 0),
	}
}

var conflictiveTypes = map[CloudflareDDNSRecordType][]CloudflareDDNSRecordType{
	RecordTypeA:     {RecordTypeCNAME},
	RecordTypeAAAA:  {RecordTypeCNAME},
	RecordTypeCNAME: {RecordTypeA, RecordTypeAAAA},
	RecordTypeTXT:   {},
}

func (c *CloudflareDDNSConfig) AddRecord(record CloudflareDDNSRecord) error {
	for _, current := range c.Records {
		if current.Name == record.Name && slices.Contains(conflictiveTypes[current.Type], record.Type) {
			return fmt.Errorf("conflictive record types for name %s: %s and %s", record.Name, current.Type, record.Type)
		}
	}

	c.Records = append(c.Records, record)
	return nil
}
