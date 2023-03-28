package golangbackend

import (
	"context"
	"errors"
	"fmt"
	errorHandler "golang-backend/internal/errors"
	"golang-backend/internal/handlers"
	"golang-backend/internal/repository"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/sirupsen/logrus"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(os.Getenv("AUTHENTICATION_SECRET")), nil)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			errorHandler.WriteError(w, errors.New("invalid token"), http.StatusUnauthorized)
			return
		}

		claimsMap := convertClaims(claims)

		userIDStr, ok := claimsMap["user_id"]
		if !ok || userIDStr == "" {
			logrus.Errorf("user id undefined")
			errorHandler.WriteError(w, errors.New("invalid token"), http.StatusUnauthorized)
			return
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			logrus.Errorf("user id undefined ParseInt")
			errorHandler.WriteError(w, errors.New("invalid token"), http.StatusUnauthorized)
			return
		}

		expStr, ok := claimsMap["exp"]
		if !ok {
			logrus.Errorf("exp undefined")
			errorHandler.WriteError(w, errors.New("token has expired"), http.StatusUnauthorized)
			return
		}

		expTime, err := time.Parse("2006-01-02 15:04:05 -0700 MST", expStr)
		if err != nil {
			logrus.Errorf("failed to parse exp claim: %v", err)
			errorHandler.WriteError(w, errors.New("invalid token"), http.StatusUnauthorized)
			return
		}

		if expTime.Unix() < time.Now().Unix() {
			logrus.Errorf("token has expired %d", userID)
			errorHandler.WriteError(w, errors.New("token has expired"), http.StatusUnauthorized)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer") {
			logrus.Errorf(err.Error())
			errorHandler.WriteError(w, errors.New("Unauthorized"), http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		tokenRepository := repository.NewTokenRepository()
		token, _ := tokenRepository.FindByToken(tokenString)
		if !token.IsActive {
			logrus.Errorf("token has expired %d", userID)
			errorHandler.WriteError(w, errors.New("invalid token"), http.StatusUnauthorized)
			return
		}

		user, err := repository.GetUserByID(userID)
		if err != nil {
			logrus.Errorf("failed to get user %d", userID)
			errorHandler.WriteError(w, errors.New("invalid token"), http.StatusInternalServerError)
			return
		}
		if user.DeletedAt != nil {
			logrus.Errorf("user is deleted %+v", user)
			errorHandler.WriteError(w, errors.New("invalid token"), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func convertClaims(claims map[string]interface{}) map[string]string {
	claimsStr := make(map[string]string)
	for key, value := range claims {
		switch v := value.(type) {
		case string:
			claimsStr[key] = v
		case int, int64:
			claimsStr[key] = strconv.FormatInt(value.(int64), 10)
		case float64:
			claimsStr[key] = strconv.FormatFloat(value.(float64), 'f', -1, 64)
		default:
			claimsStr[key] = fmt.Sprintf("%v", v)
		}
	}

	return claimsStr
}

func routes(r *chi.Mux) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.Timeout(60 * time.Second))

		r.Route("/authentication", func(r chi.Router) {
			r.Post("/login", handlers.Login)
		})

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator)
			r.Use(authMiddleware)

			r.Route("/user", func(r chi.Router) {
				r.Get("/", handlers.GetAuthenticatedUser)
				r.Post("/logout", handlers.Logout)
			})
		})
	})
}
