package main

import (
	"decoration_project/config"
	"decoration_project/routes"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// Logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config.InitDB()

	router := routes.InitializeRoutes()

	// ✅ Wrap router with CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://sajadu-c8c98.web.app"}, // your Flutter web URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(router)

	// ✅ Add logging after CORS
	loggedRouter := loggingMiddleware(corsHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8181"
	}

	log.Printf("Server started on :%s", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, loggedRouter))
}
