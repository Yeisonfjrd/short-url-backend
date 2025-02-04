package shortener

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(s *URLShortener) *gin.Engine {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:9000",
		"https://short-url-frontend-lbwn15yf6-yeisonfjrds-projects.vercel.app",
		"https://short-url-frontend-msfsxvk6v-yeisonfjrds-projects.vercel.app:3000",
		"https://short-url-frontend-git-main-yeisonfjrds-projects.vercel.app",
		"https://short-url-frontend.vercel.app",
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Accept"}
	config.AllowCredentials = false

	router.Use(cors.New(config))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := router.Group("/api")
	{
		api.POST("/shorten", s.CreateShortURL)
		api.GET("/stats/:shortCode", s.GetStats)
	}

	router.GET("/:shortCode", s.RedirectToOriginalURL)

	return router
}
