package entities

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64      `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	CreatedAt time.Time  `json:"created_at"`
	JwtTokens []JwtToken `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

func VerifyPassword(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func HashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}

func GetAuthenticatedUser(ctx context.Context) (*User, error) {
	user, ok := ctx.Value("user").(*User)
	if !ok {
		return nil, errors.New("failed to retrieve authorized user from context")
	}
	return user, nil
}
