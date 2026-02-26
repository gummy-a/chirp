package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gummy_a/chirp/media/internal/domain/value_object"
)

type contextKey string

const OwnerAccountIdKey contextKey = "OwnerAccountIdKey"

type Claims struct {
	Id domain.OwnerAccountId `json:"id"`
	jwt.RegisteredClaims
}

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		jwtSecretKey := os.Getenv("MEDIA_SERVICE_JWT_SECRET_KEY")
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// auth service signs HS256
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecretKey), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*Claims); ok {
			ctx := context.WithValue(r.Context(), OwnerAccountIdKey, claims.Id)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})
}
