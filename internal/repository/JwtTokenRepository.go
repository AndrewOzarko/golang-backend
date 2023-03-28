package repository

import (
	"errors"
	"golang-backend/internal/entities"
	"golang-backend/pkg/database"
)

type TokenRepository interface {
	CreateToken(*entities.JwtToken) error
	UpdateToken(*entities.JwtToken) error
	FindByToken(string) (*entities.JwtToken, error)
	DeleteToken(*entities.JwtToken) error
}

type tokenRepo struct{}

func NewTokenRepository() TokenRepository {
	return &tokenRepo{}
}

func (*tokenRepo) CreateToken(token *entities.JwtToken) error {
	return database.GetDB().Create(token).Error
}

func (*tokenRepo) UpdateToken(token *entities.JwtToken) error {
	return database.GetDB().Where("token = ?", token.Token).Save(token).Error
}

func (*tokenRepo) FindByToken(token string) (*entities.JwtToken, error) {
	var result entities.JwtToken
	err := database.GetDB().Where("token = ?", token).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (*tokenRepo) DeleteToken(token *entities.JwtToken) error {
	if token.Token == "" {
		return errors.New("token does not exist")
	}
	return database.GetDB().Delete(token).Error
}
