package main

import (
	"fmt"
	"log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Document struct {
	ID       string `gorm:"type:uuid;primary_key"`
	FilePath string
}

func main() {
	db, err := gorm.Open(sqlite.Open("eoffice.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var doc Document
	db.First(&doc, "id = ?", "fafc6dfa-859d-4ab0-ae04-ca1d4245253b")
	fmt.Println("FilePath:", doc.FilePath)
}
