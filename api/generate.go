package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var authAPI = []string{
	"api/auth/v1/signup.yaml",
	// "api/auth/v1/login_logout.yaml",
	// "api/auth/v1/delete_account.yaml",
}

func generateAuthServerAPI(openapi string) {
	for _, api := range authAPI {
		p := strings.Split(api, "/")
		filename := strings.Split(p[len(p)-1], ".")[0]
		outpath := "backend/servicies/auth/internal/adapter/openapi/"+filename

		currentDir, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}

		cmd := exec.Command("java", "-jar", openapi, "generate", "-i", api, "-g", "go-server", "-o", outpath, "-c", "api/config.json")

		_, err = cmd.Output()
		if err != nil {
			log.Fatal(err)
		}

		err = os.Chdir(outpath)
		if err != nil {
			log.Fatal(err)
		}

		err = os.Remove("go.mod")
		if err != nil {
			log.Fatal(err)
		}

		err = os.Remove("main.go")
		if err != nil {
			log.Fatal(err)
		}

		err = os.Chdir(currentDir)
		if err != nil {
			log.Println(err)
		}

		fmt.Printf("generated %s\n", api)
	}
}

func main() {
	openapi := os.Getenv("OPENAPI_GENERATOR_CLI_PATH")
	if(openapi == "") {
		openapi = "/usr/local/bin/openapi-generator-cli.jar"
	}

	generateAuthServerAPI(openapi)
}