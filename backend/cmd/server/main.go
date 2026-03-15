package main

import (
	"chirp/backend/router"
	authExternal "chirp/backend/services/auth/v1/external"
	mediaExternal "chirp/backend/services/media/v1/external"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf(".env not loaded.%v\n", err)
	}

	authResult, err := authExternal.CreateAuthService()
	if err != nil {
		panic(err)
	}
	defer authResult.Pool.Close()

	mediaResult, err := mediaExternal.CreateMediaService()
	if err != nil {
		panic(err)
	}
	defer mediaResult.Pool.Close()

	controller := append(authResult.Controller, mediaResult.Controller...)
	mux := router.NewAppRouter(controller)

	port := ":8080"
	fmt.Println("start server on port", port)
	http.ListenAndServe(port, mux)
}
