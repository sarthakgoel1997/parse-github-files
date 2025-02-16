package service

import (
	"context"
	"database/sql"
	"net/http"
	"parse-github-files/model"
	"time"
)

func ScanRepoJSONFiles(ctx context.Context, w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	req, err := GetScanRepoJSONFilesRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	githubReq, baseUrl, err := prepareGitHubAPIRequest(req.Repository)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rollback := true
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		http.Error(w, "error beginning db transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		if rollback {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	var fileData model.FileData

	for _, file := range req.Files {
		startTime := time.Now()

		fileData, err = getGitHubFileData(baseUrl, file, githubReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		scanResults, err := decodeAndParseBase64Data(fileData.Content)
		if err != nil {
			http.Error(w, "Failed to decode and parse: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = saveScanResults(tx, fileData.HtmlUrl, scanResults)
		if err != nil {
			http.Error(w, "Error saving scan results data to db: "+err.Error(), http.StatusInternalServerError)
			return
		}

		timeElapsed := time.Since(startTime)
		err = saveFileScannedData(tx, fileData.HtmlUrl, uint32(timeElapsed.Milliseconds()))
		if err != nil {
			http.Error(w, "Error saving file scanned data to db: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	rollback = false
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All JSON files scanned successfully"))
}
