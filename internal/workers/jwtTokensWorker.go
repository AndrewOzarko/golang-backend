package workers

import (
	"golang-backend/internal/entities"
	"golang-backend/pkg/database"
	"time"

	"github.com/sirupsen/logrus"
)

func JwtTokenWorker() {
	for {
		// Wait for 5 minutes before checking tokens again
		time.Sleep(5 * time.Minute)

		logrus.Println("JwtTokenWorker Started")

		// Get current time
		now := time.Now().Unix()

		// Start transaction
		tx := database.GetDB().Begin()

		// Get total count of tokens with is_active=true and exp < now
		var count int64
		if err := tx.Model(&entities.JwtToken{}).Where("is_active = ? AND exp < ?", true, now).Count(&count).Error; err != nil {
			tx.Rollback()
			logrus.Errorf("failed to count tokens: %v", err)
			continue
		}

		// Process tokens in chunks
		batchSize := 1000
		for offset := 0; offset < int(count); offset += batchSize {
			// Get tokens with is_active=true and exp < now
			var tokens []entities.JwtToken
			if err := tx.Where("is_active = ? AND exp < ?", true, now).Limit(batchSize).Offset(offset).Find(&tokens).Error; err != nil {
				tx.Rollback()
				logrus.Errorf("failed to find tokens: %v", err)
				continue
			}

			// Set is_active=false for expired tokens
			for _, token := range tokens {
				token.IsActive = false
				if err := tx.Model(&token).Where("token = ?", token.Token).Update("is_active", false).Error; err != nil {
					tx.Rollback()
					logrus.Errorf("failed to update token: %v", err)
					continue
				}
			}
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			logrus.Errorf("failed to commit transaction: %v", err)
			continue
		}
		logrus.Println("JwtTokenWorker Ended")
	}
}
