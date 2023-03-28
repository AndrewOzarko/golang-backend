package migrations

import (
	"golang-backend/internal/entities"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&entities.User{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entities.JwtToken{}); err != nil {
		return err
	}

	return nil
}

func Rollback(db *gorm.DB) error {
	if err := db.Migrator().DropTable(&entities.User{}); err != nil {
		return err
	}

	if err := db.Migrator().DropTable(&entities.JwtToken{}); err != nil {
		return err
	}
	return nil
}
