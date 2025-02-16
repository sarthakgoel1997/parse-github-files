package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"parse-github-files/service"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()

	// sqlite3 database configuration
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Run the SQL script to create DB tables
	err = createDbTables(db, "db.sql")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// POST API endpoint to scan GitHub repositories and save data
	router.HandleFunc("/scan", func(w http.ResponseWriter, r *http.Request) {
		service.ScanRepoJSONFiles(ctx, w, r, db)
	}).Methods("POST")

	// POST API endpoint to query stored data
	router.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		service.QueryStoredData(ctx, w, r, db)
	}).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	port := os.Getenv("PORT")
	fmt.Printf("Internal server is running on :%s...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), handler))
}

// reads and executes SQL from the file
func createDbTables(db *sql.DB, filePath string) error {
	// read the SQL file
	sqlBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading SQL file: %w", err)
	}

	// convert bytes to a string and execute it
	sqlScript := string(sqlBytes)
	_, err = db.Exec(sqlScript)
	if err != nil {
		return fmt.Errorf("error executing SQL script: %w", err)
	}

	return nil
}
