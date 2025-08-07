package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/lib/pq"
)

func main() {
	host := getEnvOrDefault("DB_HOST", "localhost")
	port := 5432
	username := getEnvOrDefault("DB_USERNAME", "admin")
	password := getEnvOrDefault("DB_PASSWORD", "1234")
	database := getEnvOrDefault("DB_NAME", "go_hexarch")
	sslMode := getEnvOrDefault("DB_SSL_MODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslMode=%s",
		host, port, username, password, database, sslMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to database successfully")

	if err := runMigrations(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migrations completed successfully")
}

func runMigrations(db *sql.DB) error {
	migrationFiles, err := filepath.Glob("*.sql")
	if err != nil {
		return fmt.Errorf("failed to find migration files: %w", err)
	}

	sort.Strings(migrationFiles)

	for _, file := range migrationFiles {
		log.Printf("Running migration: %s", file)

		content, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		log.Printf("Migration %s completed", file)
	}

	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
