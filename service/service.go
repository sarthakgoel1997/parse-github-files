package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
)

func ScanRepoJSONFiles(ctx context.Context, w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	req, err := GetScanRepoJSONFilesRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	errorMessages := processFilesAndSaveData(db, req.Repository, req.Files)
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
