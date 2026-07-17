package main

import (
	"log"

	"office-file-sharing/backend/internal/shared/config"
	"office-file-sharing/backend/internal/shared/db"
)

func main() {
	cfg := config.Load()
	gormDB := db.Init(cfg.DatabaseURL)

	log.Println("Clearing documents from the database...")

	// Use CASCADE to ensure dependent tables (workflow_histories, attachments, etc.) are truncated as well
	tx := gormDB.Exec("TRUNCATE TABLE documents CASCADE;")
	if tx.Error != nil {
		log.Fatalf("Failed to truncate documents table: %v", tx.Error)
	}

	log.Println("All documents and related records removed successfully.")
}
