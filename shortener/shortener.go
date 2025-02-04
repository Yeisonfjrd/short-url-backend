package shortener

import (
	"os"
)

type URLShortener struct {
	BaseURL string
}

func NewURLShortener() *URLShortener {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	return &URLShortener{
		BaseURL: baseURL,
	}
}
