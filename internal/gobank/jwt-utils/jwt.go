package jwt_utils

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/vortexfluc/gobank/internal/gobank/account"
	"os"
)

func ValidateJWT(tokenStr string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func CreateJWT(a *account.Account) (string, error) {
	claims := &jwt.MapClaims{
		"ExpiresAt":     15000,
		"AccountNumber": a.Number,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
