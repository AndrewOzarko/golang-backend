package handlers

import (
	"encoding/json"
	"errors"
	"golang-backend/internal/entities"
	errorHandler "golang-backend/internal/errors"
	"golang-backend/internal/repository"
	"strings"

	"net/http"

	"github.com/sirupsen/logrus"
)

func GetAuthenticatedUser(w http.ResponseWriter, r *http.Request) {
	user, err := entities.GetAuthenticatedUser(r.Context())
	if err != nil {
		logrus.Errorf(err.Error())
		errorHandler.WriteError(w, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	user, err := entities.GetAuthenticatedUser(r.Context())
	if err != nil {
		logrus.Errorf(err.Error())
		errorHandler.WriteError(w, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer") {
		logrus.Errorf(err.Error())
		errorHandler.WriteError(w, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	repository := repository.NewTokenRepository()

	token, err := repository.FindByToken(tokenString)
	if err != nil {
		logrus.Errorf(err.Error())
		errorHandler.WriteError(w, err, http.StatusUnauthorized)
		return
	}
	token.IsActive = false

	err = repository.UpdateToken(token)

	if err != nil {
		logrus.Errorf(err.Error())
		errorHandler.WriteError(w, err, http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
