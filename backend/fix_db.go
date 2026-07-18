package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=office_files port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Exec("ALTER TABLE workflow_histories ALTER COLUMN document_id DROP NOT NULL;").Error
	fmt.Println("Drop NOT NULL:", err)

	err = db.Exec("ALTER TABLE workflow_histories DROP CONSTRAINT IF EXISTS fk_workflow_histories_document;").Error
	fmt.Println("Drop FK constraint:", err)
}
