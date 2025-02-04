package shortener

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(s *URLShortener) *gin.Engine {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:9000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}
	config.AllowCredentials = true

	router.Use(cors.New(config))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	router.Static("/static", "./static")
	router.StaticFile("/", "./static/index.html")
	router.StaticFile("/app.js", "./static/app.js")

	api := router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		})
		api.POST("/shorten", s.CreateShortURL)
		api.GET("/stats/:shortCode", s.GetStats)
	}

	router.GET("/:shortCode", s.RedirectToOriginalURL)

	return router
}
