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
		panic("error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	appMux := http.NewServeMux()
	appServer := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: appMux,
	}

	repository := app.NewRepository()
	service := app.NewService()

	// registering app handlers
	app.NewHandler(appMux, repository, service)
	appMux.Handle("/static/", http.FileServer(http.Dir(".")))

	fmt.Printf("App is starting and listening on port %s...", port)
	listenErr := appServer.ListenAndServe()
	if listenErr != nil {
		panic(listenErr)
	}

}
