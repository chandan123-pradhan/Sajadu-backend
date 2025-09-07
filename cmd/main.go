package main

import (
	"decoration_project/config"
	"decoration_project/routes"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.InitDB()

	router := routes.InitializeRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8181"
	}

	log.Printf("Server started on :%s", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, router))
}
