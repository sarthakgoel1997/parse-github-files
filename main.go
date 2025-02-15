package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

func displayHelloWorld(ctx context.Context, w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Bye World!")
}

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

	// POST API endpoint to login
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		displayHelloWorld(ctx, w, r, db)
	}).Methods("GET")

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

// Reads and executes SQL from the file
func createDbTables(db *sql.DB, filePath string) error {
	// Read the SQL file
	sqlBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading SQL file: %w", err)
	}

	// Convert bytes to a string and execute it
	sqlScript := string(sqlBytes)
	_, err = db.Exec(sqlScript)
	if err != nil {
		return fmt.Errorf("error executing SQL script: %w", err)
	}

	return nil
}
