package model

type ScanRepoJSONFilesRequest struct {
	Repository string   `json:"repo"`
	Files      []string `json:"files"`
}

type QueryStoredDataRequest struct {
	Filters Filters `json:"filters"`
}

type Filters struct {
	Severity string `json:"severity"`
}

type FileData struct {
	Name     string `json:"name"`
	HtmlUrl  string `json:"html_url"`
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

type FileScanned struct {
	SourceFile string
	ScanTime   uint32
}
