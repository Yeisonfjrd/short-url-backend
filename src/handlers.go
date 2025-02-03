package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRoutes(shortener *URLShortener) *gin.Engine {
	router := gin.Default()

	router.POST("/shorten", func(c *gin.Context) {
		var req CreateURLRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		urlMapping, err := shortener.CreateShortURL(req.URL, "user", req.CustomAlias)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, URLResponse{
			ShortURL:    urlMapping.ShortCode,
			OriginalURL: urlMapping.OriginalURL,
			CreatedAt:   urlMapping.CreatedAt,
			Clicks:      urlMapping.ClickCounter,
		})
	})

	router.GET("/:shortCode", func(c *gin.Context) {
		shortCode := c.Param("shortCode")
		urlMapping, err := shortener.ResolveURL(shortCode)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, urlMapping.OriginalURL)
	})

	return router
}