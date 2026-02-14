package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gummy_a/chirp/auth/internal/domain"
)

type Claims struct {
	Id domain.AccountID `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJwt(accountID *domain.AccountID) (*string, error) {
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

func ExtractClaims(jwtToken *domain.JwtToken) (*domain.AccountID, error) {
	jwtSecretKey := os.Getenv("AUTH_SERVICE_JWT_SECRET_KEY")
    token, err := jwt.ParseWithClaims(jwtToken.String(), &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(jwtSecretKey), nil
    })

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return &claims.Id, nil
    } else {
        return nil, err
    }
}
