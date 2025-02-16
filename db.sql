CREATE TABLE IF NOT EXISTS File_Scanned (
	source_file VARCHAR(50) NOT NULL,
	scan_time INT(10) NOT NULL,
	PRIMARY KEY (source_file)
);

CREATE TABLE IF NOT EXISTS Scan_Result (
	id VARCHAR(50) NOT NULL,
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
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS Vulnerability (
	id VARCHAR(50) NOT NULL,
	scan_id VARCHAR(50) NOT NULL,
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
	PRIMARY KEY (id)
);
