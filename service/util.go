package service

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"parse-github-files/model"
	"sync"
	"time"
)

func processFilesAndSaveData(db *sql.DB, repository string, files []string) (errorMessages []string) {
	var wg sync.WaitGroup
	errChan := make(chan error, len(files)) // buffered error channel
	semaphore := make(chan struct{}, 3)     // limit concurrency to 3 files

	for _, f := range files {
		wg.Add(1)

		go func(file string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			startTime := time.Now()

			fileData, scanResults, err := getFileDataAndScanResults(repository, file)
			if err != nil {
				errChan <- fmt.Errorf("file %s: error while getting file data and scan results: %w", file, err)
				return
			}

			err = storeGitHubDataToDB(db, fileData, scanResults, startTime)
			if err != nil {
				errChan <- fmt.Errorf("file %s: error while storing GitHub data to DB: %w", file, err)
				return
			}
		}(f)
	}

	wg.Wait()      // wait for all goroutines to complete
	close(errChan) // close error channel after processing is done

	for err := range errChan {
		if err != nil {
			errorMessages = append(errorMessages, err.Error())
		}
	}
	return
}

func getFileDataAndScanResults(repo, file string) (fileData model.FileData, scanResults []model.ScanResult, err error) {
	githubReq, baseUrl, err := prepareGitHubAPIRequest(repo)
	if err != nil {
		err = fmt.Errorf("failed to prepare GitHub API request: %w", err)
		return
	}

	fileData, err = getDataFromGitHub(baseUrl, file, githubReq)
	if err != nil {
		err = fmt.Errorf("failed to fetch file data from GitHub API: %w", err)
		return
	}

	scanResults, err = decodeAndParseBase64Data(fileData.Content)
	if err != nil {
		err = fmt.Errorf("failed to decode and parse: %w", err)
	}

	return
}

func storeGitHubDataToDB(db *sql.DB, fileData model.FileData, scanResults []model.ScanResult, startTime time.Time) (err error) {
	// start a new transaction
	tx, err := db.Begin()
	if err != nil {
		err = fmt.Errorf("failed to start DB transaction: %w", err)
		return
	}
	defer tx.Rollback() // ensure rollback if anything fails

	// save scan results
	err = saveScanResults(tx, fileData.HtmlUrl, scanResults)
	if err != nil {
		err = fmt.Errorf("error saving scan results: %w", err)
		return
	}

	// save file scan metadata
	timeElapsed := time.Since(startTime)
	err = saveFileScannedData(tx, fileData.HtmlUrl, uint32(timeElapsed.Milliseconds()))
	if err != nil {
		err = fmt.Errorf("error saving file scan metadata: %w", err)
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		err = fmt.Errorf("failed to commit transaction: %w", err)
	}
	return
}

func decodeAndParseBase64Data(encoded string) ([]model.ScanResult, error) {
	// decode Base64
	jsonData, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("error decoding base64: %v", err)
	}

	// unmarshal JSON into struct array
	var scanResultsWrappers []model.ScanResultsWrapper
	err = json.Unmarshal(jsonData, &scanResultsWrappers)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	// extract ScanResults from wrapper
	var scanResults []model.ScanResult
	for _, wrapper := range scanResultsWrappers {
		scanResults = append(scanResults, wrapper.ScanResults)
	}

	return scanResults, nil
}

func convertArrayToJson(req []string) string {
	jsonData, _ := json.Marshal(req)
	return string(jsonData)
}
