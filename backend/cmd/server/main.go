package main

import (
	"chirp/backend/router"
	authExternal "chirp/backend/services/auth/v1/external"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf(".env not loaded.%v\n", err)
	}

	authResult := authExternal.CreateAuthService()
	defer authResult.Pool.Close()
	mux := router.NewAppRouter(authResult.Controller)

	port := ":8080"
	fmt.Println("start server on port", port)
	http.ListenAndServe(port, mux)
}
