package database

import (
	"fmt"
	"log"

	"golang-backend/internal/migrations"
	"golang-backend/pkg/database"
)

func Migrate() {
	db := database.GetDB()
	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate: %s", err)
	}
	fmt.Println("Migration successful")
}

func Rollback() {
	db := database.GetDB()
	if err := migrations.Rollback(db); err != nil {
		log.Fatalf("Failed to rollback: %s", err)
	}
	fmt.Println("Rollback successful")
}
