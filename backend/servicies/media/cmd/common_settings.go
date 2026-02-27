package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func SetDefaultEnvironmentVariables() {
	os.Setenv("MEDIA_SERVICE_APP_ENV", "development")
	env := os.Getenv("MEDIA_SERVICE_APP_ENV")
	if env == "development" {
		os.Setenv("MEDIA_SERVICE_PORT", "8081")
		os.Setenv("MEDIA_SERVICE_JWT_SECRET_KEY", "PSsDWRYMnGnLZpq1uq4Dd24WnGncTBkbtciiXzFNqGPHyJ") // must be same as auth service jwt-secret-key
		os.Setenv("MEDIA_SERVICE_ALLOW_ORIGIN", "http://localhost:3000")                            // DO NOT SET WILDCARD
		os.Setenv("MEDIA_SERVICE_DATABASE_URL", "postgres://postgres:password@localhost:5432/media_service?sslmode=disable")
		os.Setenv("MEDIA_SERVICE_REDIS_URL", "localhost:6379")
	} else {
		err := godotenv.Load()
		if err != nil {
			fmt.Printf(".env not loaded. %v\n", err)
		}
	}
}

func CheckEnvironmentVariables() {
	env := os.Getenv("MEDIA_SERVICE_APP_ENV")
	if env == "" {
		log.Fatal("MEDIA_SERVICE_APP_ENV environment variable is not set")
	}

	port := os.Getenv("MEDIA_SERVICE_PORT")
	if env == "production" && port == "" {
		log.Fatal("MEDIA_SERVICE_PORT environment variable is required in production")
	}

	jwtSecretKey := os.Getenv("MEDIA_SERVICE_JWT_SECRET_KEY")
	if env == "production" && jwtSecretKey == "" {
		log.Fatal("MEDIA_SERVICE_JWT_SECRET_KEY is not set")
	}

	allowOrigin := os.Getenv("MEDIA_SERVICE_ALLOW_ORIGIN")
	if env == "production" && allowOrigin == "" {
		log.Fatal("MEDIA_SERVICE_ALLOW_ORIGIN environment variable is required in production")
	}

	url := os.Getenv("MEDIA_SERVICE_DATABASE_URL")
	if env == "production" && url == "" {
		log.Fatal("MEDIA_SERVICE_DATABASE_URL is not set")
	}

	redis := os.Getenv("MEDIA_SERVICE_REDIS_URL")
	if env == "production" && redis == "" {
		log.Fatal("MEDIA_SERVICE_REDIS_URL is not set")
	}
}
