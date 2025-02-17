package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

func ScanRepoJSONFiles(ctx context.Context, w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	req, err := GetScanRepoJSONFilesRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(req.Files)) // buffered error channel
	semaphore := make(chan struct{}, 3)         // limit concurrency to 3 files

	for _, f := range req.Files {
		wg.Add(1)

		go func(file string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			startTime := time.Now()

			fileData, scanResults, err := getFileDataAndScanResults(req.Repository, file)
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

	var errorMessages []string // handle errors
	for err := range errChan {
		if err != nil {
			errorMessages = append(errorMessages, err.Error())
		}
	}
	if len(errorMessages) > 0 {
		http.Error(w, strings.Join(errorMessages, "\n"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All JSON files scanned successfully"))
}

func QueryStoredData(ctx context.Context, w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	req, err := QueryStoredDataRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := getFilteredData(db, req)
	if err != nil {
		http.Error(w, "error while getting data from DB: "+err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(resp)
}
