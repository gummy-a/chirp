package jwt

import (
	"chirp/backend/services/auth/v1/internal/domain/value_object"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Id value_object.AccountID `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJwt(accountID *value_object.AccountID) (*string, error) {
	claims := Claims{
		Id: *accountID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "app",
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecretKey := os.Getenv("AUTH_SERVICE_JWT_SECRET_KEY")
	jwtResult, err := jwtToken.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return nil, err
	}
	return &jwtResult, nil
}
