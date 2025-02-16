package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"parse-github-files/model"
)

func GetScanRepoJSONFilesRequest(w http.ResponseWriter, r *http.Request) (req model.ScanRepoJSONFilesRequest, err error) {
	// parse the incoming JSON data from the request body
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return req, fmt.Errorf("invalid request body: %v", err)
	}

	// validate the request
	if req.Repository == "" {
		return req, fmt.Errorf("GitHub repository cannot be empty")
	}
	if len(req.Files) == 0 {
		return req, fmt.Errorf("no files to scan")
	}

	return
}
