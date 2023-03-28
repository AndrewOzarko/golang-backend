package jwt

import (
	"log"
	"os"

	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	secret := os.Getenv("AUTHENTICATION_SECRET")
	if secret == "" {
		log.Fatalln("Need to generate AUTHENTICATION_SECRET before use")
	}
	tokenAuth = jwtauth.New("HS256", []byte(secret), nil)
}

func GenerateJwtString(claims map[string]interface{}) (string, error) {
	_, token, err := tokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}
