package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"go_test_task_4/internal/app"
	"net/http"
	"os"
)

func main() {
	envLoadErr := godotenv.Load(".env")
	if envLoadErr != nil {
		panic("Error loading .env file")
	}
	port := os.Getenv("PORT")

	appMux := http.NewServeMux()
	appServer := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: appMux,
	}

	// registering app handlers
	app.NewHandler(appMux)

	fmt.Printf("App is starting and listening on port %s...", port)
	listenErr := appServer.ListenAndServe()
	if listenErr != nil {
		return
	}

}
