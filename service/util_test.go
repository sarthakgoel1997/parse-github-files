package service

import (
	"database/sql"
	"parse-github-files/model"
	"reflect"
	"testing"
	"time"
)

func Test_convertArrayToJson(t *testing.T) {
	type args struct {
		req []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "array with 0 elements",
			args: args{
				req: []string{},
			},
			want: `[]`,
		},
		{
			name: "array with 1 element",
			args: args{
				req: []string{"test"},
			},
			want: `["test"]`,
		},
		{
			name: "array with 2 elements",
			args: args{
				req: []string{"test1", "test2"},
			},
			want: `["test1","test2"]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertArrayToJson(tt.args.req); got != tt.want {
				t.Errorf("convertArrayToJson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeAndParseBase64Data(t *testing.T) {
	type args struct {
		encoded string
	}
	tests := []struct {
		name    string
		args    args
		want    []model.ScanResult
		wantErr bool
	}{
		{
			name: "decoding of correctly encoded string",
			args: args{
				encoded: "WwogIHsKICAgICJzY2FuUmVzdWx0cyI6IHsKICAgICAgInNjYW5faWQiOiAiVlVMTl9zY2FuXzM0NW1ubyIsCiAgICAgICJ0aW1lc3RhbXAiOiAiMjAyNS0wMS0yOVQxMzowMDowMFoiLAogICAgICAic2Nhbl9zdGF0dXMiOiAiY29tcGxldGVkIiwKICAgICAgInJlc291cmNlX3R5cGUiOiAiY29udGFpbmVyIiwKICAgICAgInJlc291cmNlX25hbWUiOiAibWwtaW5mZXJlbmNlOjIuMC4wIiwKICAgICAgInZ1bG5lcmFiaWxpdGllcyI6IFsKICAgICAgICB7CiAgICAgICAgICAiaWQiOiAiQ1ZFLTIwMjQtNTU1NSIsCiAgICAgICAgICAic2V2ZXJpdHkiOiAiSElHSCIsCiAgICAgICAgICAiY3ZzcyI6IDguNSwKICAgICAgICAgICJzdGF0dXMiOiAiYWN0aXZlIiwKICAgICAgICAgICJwYWNrYWdlX25hbWUiOiAidGVuc29yZmxvdyIsCiAgICAgICAgICAiY3VycmVudF92ZXJzaW9uIjogIjIuNy4wIiwKICAgICAgICAgICJmaXhlZF92ZXJzaW9uIjogIjIuNy4xIiwKICAgICAgICAgICJkZXNjcmlwdGlvbiI6ICJSZW1vdGUgY29kZSBleGVjdXRpb24gaW4gVGVuc29yRmxvdyBtb2RlbCBsb2FkaW5nIiwKICAgICAgICAgICJwdWJsaXNoZWRfZGF0ZSI6ICIyMDI1LTAxLTI0VDAwOjAwOjAwWiIsCiAgICAgICAgICAibGluayI6ICJodHRwczovL252ZC5uaXN0Lmdvdi92dWxuL2RldGFpbC9DVkUtMjAyNC01NTU1IiwKICAgICAgICAgICJyaXNrX2ZhY3RvcnMiOiBbCiAgICAgICAgICAgICJSZW1vdGUgQ29kZSBFeGVjdXRpb24iLAogICAgICAgICAgICAiSGlnaCBDVlNTIFNjb3JlIiwKICAgICAgICAgICAgIlB1YmxpYyBFeHBsb2l0IEF2YWlsYWJsZSIsCiAgICAgICAgICAgICJFeHBsb2l0IGluIFdpbGQiCiAgICAgICAgICBdCiAgICAgICAgfQogICAgICBdLAogICAgICAic3VtbWFyeSI6IHsKICAgICAgICAidG90YWxfdnVsbmVyYWJpbGl0aWVzIjogMywKICAgICAgICAic2V2ZXJpdHlfY291bnRzIjogewogICAgICAgICAgIkNSSVRJQ0FMIjogMCwKICAgICAgICAgICJISUdIIjogMSwKICAgICAgICAgICJNRURJVU0iOiAxLAogICAgICAgICAgIkxPVyI6IDEKICAgICAgICB9LAogICAgICAgICJmaXhhYmxlX2NvdW50IjogMywKICAgICAgICAiY29tcGxpYW50IjogZmFsc2UKICAgICAgfSwKICAgICAgInNjYW5fbWV0YWRhdGEiOiB7CiAgICAgICAgInNjYW5uZXJfdmVyc2lvbiI6ICIzMC4xLjUxIiwKICAgICAgICAicG9saWNpZXNfdmVyc2lvbiI6ICIyMDI1LjEuMjkiLAogICAgICAgICJzY2FubmluZ19ydWxlcyI6IFsKICAgICAgICAgICJ2dWxuZXJhYmlsaXR5IiwKICAgICAgICAgICJjb21wbGlhbmNlIiwKICAgICAgICAgICJtYWx3YXJlIgogICAgICAgIF0sCiAgICAgICAgImV4Y2x1ZGVkX3BhdGhzIjogWwogICAgICAgICAgIi90bXAiLAogICAgICAgICAgIi92YXIvbG9nIgogICAgICAgIF0KICAgICAgfQogICAgfQogIH0KXQ==",
			},
			want: []model.ScanResult{
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
			},
			wantErr: false,
		},
		{
			name: "error decoding encoded string",
			args: args{
				encoded: "ABCD",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeAndParseBase64Data(tt.args.encoded)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeAndParseBase64Data() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeAndParseBase64Data() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_processFilesAndSaveData(t *testing.T) {
	type args struct {
		db         *sql.DB
		repository string
		files      []string
	}
	tests := []struct {
		name              string
		args              args
		wantErrorMessages []string
	}{
		{
			name: "success call without any error messages",
			args: args{
				db: func() *sql.DB {
					db, _ := sql.Open("sqlite3", ":memory:")
					_, _ = db.Exec(`
								CREATE TABLE IF NOT EXISTS File_Scanned (
									source_file VARCHAR(50) NOT NULL,
									scan_time INT(10) NOT NULL,
									PRIMARY KEY (source_file)
								);

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
								);`)
					return db
				}(),
				repository: "https://github.com/velancio/vulnerability_scans",
				files:      []string{"vulnscan1011.json", "vulnscan1213.json"},
			},
			wantErrorMessages: nil,
		},
		{
			name: "file does not exist in the GitHub repo",
			args: args{
				db: func() *sql.DB {
					db, _ := sql.Open("sqlite3", ":memory:")
					_, _ = db.Exec(`
								CREATE TABLE IF NOT EXISTS File_Scanned (
									source_file VARCHAR(50) NOT NULL,
									scan_time INT(10) NOT NULL,
									PRIMARY KEY (source_file)
								);

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
								);`)
					return db
				}(),
				repository: "https://github.com/velancio/vulnerability_scans",
				files:      []string{"abcd.json"},
			},
			wantErrorMessages: []string{"file abcd.json: error while getting file data and scan results: failed to fetch file data from GitHub API: GitHub API request failed with status code: 404"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErrorMessages := processFilesAndSaveData(tt.args.db, tt.args.repository, tt.args.files); !reflect.DeepEqual(gotErrorMessages, tt.wantErrorMessages) {
				t.Errorf("processFilesAndSaveData() = %v, want %v", gotErrorMessages, tt.wantErrorMessages)
			}
		})
	}
}
