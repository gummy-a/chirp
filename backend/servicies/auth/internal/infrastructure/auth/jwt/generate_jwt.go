package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gummy_a/chirp/auth/internal/domain"
)

type Claims struct {
	Id domain.AccountID `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJwt(accountID domain.AccountID) (*string, error) {

	claims := Claims{
		Id: accountID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "app",
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecretKey := os.Getenv("AUTH_SERVICE_JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		return nil, errors.New("AUTH_SERVICE_JWT_SECRET_KEY is not set")
	}

	jwtResult, err := jwtToken.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return nil, err
	}
	return &jwtResult, nil
}
