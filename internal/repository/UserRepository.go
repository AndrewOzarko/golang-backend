package repository

import (
	"errors"
	"golang-backend/internal/entities"
	"golang-backend/pkg/database"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func GetUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	result := database.GetDB().Where(&entities.User{Email: email}).Where("deleted_at IS NULL").First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetUserByID(id int64) (*entities.User, error) {
	var user entities.User
	result := database.GetDB().Where(&entities.User{ID: id}).Where("deleted_at IS NULL").First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logrus.Errorf(result.Error.Error())
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &user, nil
}
