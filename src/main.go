package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"url-shortener/database"
	"url-shortener/shortener"
)

func main() {
	_ = godotenv.Load()

	database.InitDatabase()

	urlShortener := shortener.NewURLShortener()
	router := shortener.SetupRoutes(urlShortener)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	router.Run(":" + port)
}
