package service

import (
	"database/sql"
	"encoding/json"
	"parse-github-files/model"
)

func saveScanResults(tx *sql.Tx, sourceFile string, scanResults []model.ScanResult) (err error) {
	var query string

	for _, res := range scanResults {
		// add scan result
		query = queryToAddScanResult()
		_, err = tx.Exec(query, res.ID, sourceFile, res.Timestamp, res.ScanStatus, res.ResourceType, res.ResourceName, res.Summary.TotalVulnerabilities, res.Summary.SeverityCounts.Critical, res.Summary.SeverityCounts.High, res.Summary.SeverityCounts.Medium, res.Summary.SeverityCounts.Low, res.Summary.FixableCount, res.Summary.Compliant, res.ScanMetadata.ScannerVersion, res.ScanMetadata.PoliciesVersion, convertArrayToJson(res.ScanMetadata.ScanningRules), convertArrayToJson(res.ScanMetadata.ExcludedPaths))
		if err != nil {
			return
		}

		// add related vulnerabilities
		for _, v := range res.Vulnerabilities {
			query = queryToAddVulnerability()
			_, err = tx.Exec(query, v.ID, res.ID, v.Severity, v.Cvss, v.Status, v.PackageName, v.CurrentVersion, v.FixedVersion, v.Description, v.PublishedDate, v.Link, convertArrayToJson(v.RiskFactors))
			if err != nil {
				return
			}
		}
	}
	return
}

func getFilteredData(db *sql.DB, req model.QueryStoredDataRequest) (resp []model.Vulnerability, err error) {
	resp = []model.Vulnerability{}
	query := queryToGetFilteredVulnerability()
	rows, err := db.Query(query, req.Filters.Severity)
	if err != nil {
		return
	}
	defer rows.Close()

	var v model.Vulnerability
	var riskFactors string

	for rows.Next() {
		err = rows.Scan(&v.ID, &v.ScanID, &v.Severity, &v.Cvss, &v.Status, &v.PackageName, &v.CurrentVersion, &v.FixedVersion, &v.Description, &v.PublishedDate, &v.Link, &riskFactors)
		if err != nil {
			return
		}
		json.Unmarshal([]byte(riskFactors), &v.RiskFactors)
		resp = append(resp, v)
	}
	return
}

func saveFileScannedData(tx *sql.Tx, sourceFile string, scanTime uint32) (err error) {
	query := queryToAddFileScannedData()
	_, err = tx.Exec(query, sourceFile, scanTime)
	return
}

func queryToAddFileScannedData() string {
	sqlQuery := `
				INSERT INTO File_Scanned
					(source_file, scan_time)
				VALUES
					(?, ?);
				`
	return sqlQuery
}

func queryToAddScanResult() string {
	sqlQuery := `
				INSERT INTO Scan_Result
					(id, source_file, timestamp, scan_status, resource_type, resource_name, total_vulnerabilities, critical_severity, high_severity, medium_severity, low_severity, fixable_count, compliant, scanner_version, policies_version, scanning_rules, excluded_paths)
				VALUES
					(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
				`
	return sqlQuery
}

func queryToAddVulnerability() string {
	sqlQuery := `
				INSERT INTO Vulnerability
					(id, scan_id, severity, cvss, status, package_name, current_version, fixed_version, description, published_date, link, risk_factors)
				VALUES
					(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
				`
	return sqlQuery
}

func queryToGetFilteredVulnerability() string {
	sqlQuery := `
	SELECT
		*
	FROM
		Vulnerability
	WHERE
		severity = ?;
	`
	return sqlQuery
}
