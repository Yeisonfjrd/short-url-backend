package main

import (
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type URLShortener struct {
	redis *redis.Client
}

func NewURLShortener(redisClient *redis.Client) *URLShortener {
	return &URLShortener{redis: redisClient}
}

func generateShortCode(url string) string {
	hash := sha256.Sum256([]byte(url + uuid.New().String()))
	shortCode := base64.URLEncoding.EncodeToString(hash[:])
	return strings.TrimRight(shortCode[:8], "=")
}

func (us *URLShortener) CreateShortURL(originalURL, createdBy string, customAlias string) (*URLMapping, error) {
	var shortCode string
	if customAlias != "" {
		shortCode = customAlias
	} else {
		shortCode = generateShortCode(originalURL)
	}

	urlMapping := &URLMapping{
		ShortCode:     shortCode,
		OriginalURL:   originalURL,
		CreatedBy:     createdBy,
		Enabled:       true,
		ClickCounter:  0,
		LastAccessed:  time.Now(),
	}

	result := DB.Create(urlMapping)
	if result.Error != nil {
		return nil, result.Error
	}

	return urlMapping, nil
}

func (us *URLShortener) ResolveURL(shortCode string) (*URLMapping, error) {
	var urlMapping URLMapping
	result := DB.Where("short_code = ? AND enabled = true", shortCode).First(&urlMapping)
	if result.Error != nil {
		return nil, result.Error
	}

	urlMapping.LastAccessed = time.Now()
	urlMapping.ClickCounter++
	DB.Save(&urlMapping)

	return &urlMapping, nil
}