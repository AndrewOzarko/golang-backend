package seeders

import (
	"golang-backend/internal/entities"

	"gorm.io/gorm"
)

func CreateUsers(db *gorm.DB) error {
	users := []entities.User{
		{Email: "admin@admin.cc", Password: entities.HashPassword("secret")},
		{Email: "jane@example.com", Password: entities.HashPassword("secret")},
	}

	for _, u := range users {
		if err := db.Create(&u).Error; err != nil {
			return err
		}
	}

	return nil
}
