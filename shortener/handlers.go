package shortener

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"time"

	"url-shortener/database"

	"github.com/gin-gonic/gin"
)

type ShortenRequest struct {
	URL string `json:"url" binding:"required"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

func (s *URLShortener) CreateShortURL(c *gin.Context) {
	var req ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL is required"})
		return
	}

	hash := sha256.Sum256([]byte(req.URL))
	shortCode := base64.URLEncoding.EncodeToString(hash[:8])

	var existingURL database.URL
	result := database.GetDatabase().Where("original_url = ?", req.URL).First(&existingURL)
	if result.Error == nil {
		shortURL := s.BaseURL + "/" + existingURL.ShortCode
		c.JSON(http.StatusOK, ShortenResponse{ShortURL: shortURL})
		return
	}

	url := database.URL{
		OriginalURL: req.URL,
		ShortCode:   shortCode,
	}

	result = database.GetDatabase().Create(&url)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store URL"})
		return
	}

	shortURL := s.BaseURL + "/" + shortCode
	c.JSON(http.StatusOK, ShortenResponse{ShortURL: shortURL})
}

func (s *URLShortener) RedirectToOriginalURL(c *gin.Context) {
	shortCode := c.Param("shortCode")

	var url database.URL
	result := database.GetDatabase().Where("short_code = ?", shortCode).First(&url)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	database.GetDatabase().Model(&url).Updates(database.URL{
		Visits:    url.Visits + 1,
		LastVisit: time.Now(),
	})

	c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
}

func (s *URLShortener) GetStats(c *gin.Context) {
	shortCode := c.Param("shortCode")

	var url database.URL
	result := database.GetDatabase().Where("short_code = ?", shortCode).First(&url)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short_code": shortCode,
		"url":        url.OriginalURL,
		"visits":     url.Visits,
		"last_visit": url.LastVisit,
		"created_at": url.CreatedAt,
	})
}
