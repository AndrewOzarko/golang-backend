package handlers

import (
	"encoding/json"
	"errors"
	"golang-backend/internal/entities"
	errorHandler "golang-backend/internal/errors"
	"golang-backend/internal/repository"
	"golang-backend/pkg/jwt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var validate = validator.New()

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5"`
}

type LoginResponse struct {
	ID          int64  `json:"id"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorHandler.WriteError(w, err, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		errorHandler.WriteError(w, err, http.StatusBadRequest)
		return
	}

	user, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorHandler.WriteError(w, errors.New("invalid email or password"), http.StatusUnauthorized)
		} else {
			logrus.Errorf("Failed to get user by email: %v", err)
			errorHandler.WriteError(w, errors.New("failed to get user"), http.StatusInternalServerError)
		}
		return
	}

	if err := entities.VerifyPassword(req.Password, user.Password); err != nil {
		errorHandler.WriteError(w, errors.New("invalid email or password"), http.StatusUnauthorized)
		return
	}

	exp := time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

	claims := map[string]interface{}{
		"user_id": user.ID,
		"exp":     exp,
	}
	tokenString, err := jwt.GenerateJwtString(claims)
	if err != nil {
		logrus.Errorf("Failed to generate jwt string: %v", err)
		errorHandler.WriteError(w, errors.New("failed to generate access token"), http.StatusInternalServerError)
		return
	}

	device := r.Header.Get("User-Agent")

	repository := repository.NewTokenRepository()

	repository.CreateToken(&entities.JwtToken{
		UserID: user.ID,
		Exp:    exp,
		Device: device,
		Token:  tokenString,
	})

	resp := LoginResponse{
		ID:          user.ID,
		Email:       user.Email,
		AccessToken: tokenString,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
