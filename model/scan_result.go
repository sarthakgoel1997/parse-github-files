package model

import "time"

type ScanResultsWrapper struct {
	ScanResults ScanResult `json:"scanResults"`
}

type ScanResult struct {
	ScanID          string          `json:"scan_id"`
	SourceFile      string          `json:"source_file"`
	Timestamp       time.Time       `json:"timestamp"`
	ScanStatus      string          `json:"scan_status"`
	ResourceType    string          `json:"resource_type"`
	ResourceName    string          `json:"resource_name"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	Summary         Summary         `json:"summary"`
	ScanMetadata    Metadata        `json:"scan_metadata"`
}

type Summary struct {
	TotalVulnerabilities uint32        `json:"total_vulnerabilities"`
	SeverityCounts       SeverityCount `json:"severity_counts"`
	FixableCount         uint32        `json:"fixable_count"`
	Compliant            bool          `json:"compliant"`
}

type SeverityCount struct {
	Critical uint32 `json:"CRITICAL"`
	High     uint32 `json:"HIGH"`
	Medium   uint32 `json:"MEDIUM"`
	Low      uint32 `json:"LOW"`
}

type Metadata struct {
	ScannerVersion  string   `json:"scanner_version"`
	PoliciesVersion string   `json:"policies_version"`
	ScanningRules   []string `json:"scanning_rules"`
	ExcludedPaths   []string `json:"excluded_paths"`
}
