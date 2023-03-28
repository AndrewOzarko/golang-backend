package database

import (
	"fmt"
	"golang-backend/internal/seeders"
	"golang-backend/pkg/database"
	"log"
)

func Seed() {
	db := database.GetDB()
	if err := seeders.CreateUsers(db); err != nil {
		log.Fatalf("Failed to seed: %s", err)
	}
	fmt.Println("Seeding successful")
}
