package service

import (
	"database/sql"
	"encoding/json"
	"parse-github-files/model"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestSaveScanResults_Success(t *testing.T) {
	// create in-memory SQLite DB
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	// create mock schema
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS Scan_Result (
			scan_id VARCHAR(50) NOT NULL,
			source_file VARCHAR(50) NOT NULL,
			timestamp TIMESTAMP NOT NULL,
			scan_status VARCHAR(30) NOT NULL,
			resource_type VARCHAR(30) NOT NULL,
			resource_name VARCHAR(50) NOT NULL,
			total_vulnerabilities INT(10) NOT NULL,
			critical_severity INT(10) NOT NULL,
			high_severity INT(10) NOT NULL,
			medium_severity INT(10) NOT NULL,
			low_severity INT(10) NOT NULL,
			fixable_count INT(10) NOT NULL,
			compliant BOOLEAN NOT NULL,
			scanner_version VARCHAR(30) NOT NULL,
			policies_version VARCHAR(30) NOT NULL,
			scanning_rules TEXT NOT NULL,
			excluded_paths TEXT NOT NULL,
			PRIMARY KEY (scan_id, source_file)
		);
		CREATE TABLE IF NOT EXISTS Vulnerability (
			id VARCHAR(50) NOT NULL,
			scan_id VARCHAR(50) NOT NULL,
			source_file VARCHAR(50) NOT NULL,
			severity VARCHAR(20) NOT NULL,
			cvss DOUBLE NOT NULL,
			status VARCHAR(20) NOT NULL,
			package_name VARCHAR(50) NOT NULL,
			current_version VARCHAR(20) NOT NULL,
			fixed_version VARCHAR(20) NOT NULL,
			description TEXT NOT NULL,
			published_date TIMESTAMP NOT NULL,
			link VARCHAR(200) NOT NULL,
			risk_factors TEXT NOT NULL,
			PRIMARY KEY (id, scan_id, source_file)
		);
	`)
	assert.NoError(t, err)

	// begin transaction
	tx, err := db.Begin()
	assert.NoError(t, err)

	// test data
	scanResults := []model.ScanResult{
		{
			ScanID: "VULN_scan_345mno",
			Timestamp: func() time.Time {
				t, _ := time.Parse(time.RFC3339, "2025-01-29T13:00:00Z")
				return t
			}(),
			ScanStatus:   "completed",
			ResourceType: "container",
			ResourceName: "ml-inference:2.0.0",
			Vulnerabilities: []model.Vulnerability{
				{
					ID:             "CVE-2024-5555",
					Severity:       "HIGH",
					Cvss:           8.5,
					Status:         "active",
					PackageName:    "tensorflow",
					CurrentVersion: "2.7.0",
					FixedVersion:   "2.7.1",
					Description:    "Remote code execution in TensorFlow model loading",
					PublishedDate: func() time.Time {
						t, _ := time.Parse(time.RFC3339, "2025-01-24T00:00:00Z")
						return t
					}(),
					Link:        "https://nvd.nist.gov/vuln/detail/CVE-2024-5555",
					RiskFactors: []string{"Remote Code Execution", "High CVSS Score", "Public Exploit Available", "Exploit in Wild"},
				},
			},
			Summary: model.Summary{
				TotalVulnerabilities: 3,
				SeverityCounts: model.SeverityCount{
					Critical: 0,
					High:     1,
					Medium:   1,
					Low:      1,
				},
				FixableCount: 3,
				Compliant:    false,
			},
			ScanMetadata: model.Metadata{
				ScannerVersion:  "30.1.51",
				PoliciesVersion: "2025.1.29",
				ScanningRules:   []string{"vulnerability", "compliance", "malware"},
				ExcludedPaths:   []string{"/tmp", "/var/log"},
			},
		},
	}

	// call function
	err = saveScanResults(tx, "source.json", scanResults)
	assert.NoError(t, err)

	// verify scan_results table
	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM Scan_Result WHERE scan_id = ?", "VULN_scan_345mno").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count, "Scan result should be inserted")

	// verify vulnerabilities table
	err = tx.QueryRow("SELECT COUNT(*) FROM Vulnerability WHERE id = ?", "CVE-2024-5555").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count, "Vulnerability should be inserted")

	// rollback transaction to keep DB clean
	tx.Rollback()
}

func TestSaveScanResults_SQLFailure(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	// create only Vulnerability table (skip Scan_Result to induce failure)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS Vulnerability (
			id VARCHAR(50) NOT NULL,
			scan_id VARCHAR(50) NOT NULL,
			source_file VARCHAR(50) NOT NULL,
			severity VARCHAR(20) NOT NULL,
			cvss DOUBLE NOT NULL,
			status VARCHAR(20) NOT NULL,
			package_name VARCHAR(50) NOT NULL,
			current_version VARCHAR(20) NOT NULL,
			fixed_version VARCHAR(20) NOT NULL,
			description TEXT NOT NULL,
			published_date TIMESTAMP NOT NULL,
			link VARCHAR(200) NOT NULL,
			risk_factors TEXT NOT NULL,
			PRIMARY KEY (id, scan_id, source_file)
		);
	`)
	assert.NoError(t, err)

	tx, err := db.Begin()
	assert.NoError(t, err)

	// test data with missing scan_results table
	scanResults := []model.ScanResult{
		{
			ScanID:       "test_scan",
			Timestamp:    time.Now(),
			ScanStatus:   "completed",
			ResourceType: "container",
			ResourceName: "test-container",
		},
	}

	// call function (should fail)
	err = saveScanResults(tx, "source.json", scanResults)
	assert.Error(t, err)

	tx.Rollback()
}

func TestSaveFileScannedData_Success(t *testing.T) {
	// create in-memory SQLite DB
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	// create mock schema
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS File_Scanned (
		source_file VARCHAR(50) NOT NULL,
		scan_time INT(10) NOT NULL,
		PRIMARY KEY (source_file)
		);
	`)
	assert.NoError(t, err)

	// begin transaction
	tx, err := db.Begin()
	assert.NoError(t, err)

	// call function
	err = saveFileScannedData(tx, "vulnscan1011.json", 200)
	assert.NoError(t, err)

	// verify File_Scanned table
	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM File_Scanned WHERE source_file = ?", "vulnscan1011.json").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count, "File_Scanned result should be inserted")

	err = tx.QueryRow("SELECT COUNT(*) FROM File_Scanned WHERE source_file = ?", "abcd.json").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count, "Any other random source file does not exist")

	// rollback transaction to keep DB clean
	tx.Rollback()
}

func TestGetFilteredData(t *testing.T) {
	// setup in-memory SQLite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory SQLite: %v", err)
	}
	defer db.Close()

	// create tabl for testing
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS Vulnerability (
		id VARCHAR(50) NOT NULL,
		scan_id VARCHAR(50) NOT NULL,
		source_file VARCHAR(50) NOT NULL,
		severity VARCHAR(20) NOT NULL,
		cvss DOUBLE NOT NULL,
		status VARCHAR(20) NOT NULL,
		package_name VARCHAR(50) NOT NULL,
		current_version VARCHAR(20) NOT NULL,
		fixed_version VARCHAR(20) NOT NULL,
		description TEXT NOT NULL,
		published_date TIMESTAMP NOT NULL,
		link VARCHAR(200) NOT NULL,
		risk_factors TEXT NOT NULL,
		PRIMARY KEY (id, scan_id, source_file)
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// add mock data
	mockVulnerabilities := []model.Vulnerability{
		{
			ID:             "CVE-2024-2224",
			ScanID:         "VULN_scan_458def",
			SourceFile:     "https://github.com/velancio/vulnerability_scans/blob/main/vulnscan1011.json",
			Severity:       "MEDIUM",
			Cvss:           8.8,
			Status:         "completed",
			PackageName:    "container",
			CurrentVersion: "5.6.2",
			FixedVersion:   "5.6.6",
			Description:    "Authentication bypass in Spring Security",
			PublishedDate:  time.Now(),
			Link:           "https://nvd.nist.gov/vuln/detail/CVE-2024-2222",
			RiskFactors:    []string{"Authentication Bypass", "High CVSS Score", "Proof of Concept Exploit Available"},
		},
		{
			ID:             "CVE-2024-2222",
			ScanID:         "VULN_scan_456def",
			SourceFile:     "https://github.com/velancio/vulnerability_scans/blob/main/vulnscan1011.json",
			Severity:       "HIGH",
			Cvss:           8.2,
			Status:         "active",
			PackageName:    "spring-security",
			CurrentVersion: "5.6.0",
			FixedVersion:   "5.6.1",
			Description:    "Authentication bypass in Spring Security",
			PublishedDate:  time.Now(),
			Link:           "https://nvd.nist.gov/vuln/detail/CVE-2024-2222",
			RiskFactors:    []string{"Authentication Bypass", "High CVSS Score", "Proof of Concept Exploit Available"},
		},
	}

	for _, v := range mockVulnerabilities {
		riskFactorsJSON, _ := json.Marshal(v.RiskFactors)
		_, err = db.Exec(
			`INSERT INTO Vulnerability (id, scan_id, source_file, severity, cvss, status, package_name, current_version, fixed_version, description, published_date, link, risk_factors) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			v.ID, v.ScanID, v.SourceFile, v.Severity, v.Cvss, v.Status, v.PackageName, v.CurrentVersion, v.FixedVersion, v.Description, v.PublishedDate, v.Link, string(riskFactorsJSON),
		)
		if err != nil {
			t.Fatalf("Failed to insert mock data: %v", err)
		}
	}

	// call function with HIGH severity filter
	req := model.QueryStoredDataRequest{
		Filters: model.Filters{Severity: "HIGH"},
	}
	result, err := getFilteredData(db, req)
	if err != nil {
		t.Fatalf("Error executing getFilteredData: %v", err)
	}

	// validate expected results
	expectedCount := 1
	if len(result) != expectedCount {
		t.Errorf("Expected %d results, got %d", expectedCount, len(result))
	}

	if result[0].ID != "CVE-2024-2222" || result[0].Severity != "HIGH" {
		t.Errorf("Unexpected result: %+v", result[0])
	}
}
