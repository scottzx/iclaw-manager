package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"iclaw-admin-api/internal/router"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "18789"
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Setup API routes
	router.SetupRouter(r)

	log.Printf("iClaw Admin API starting on port %s", port)
	r.Run("0.0.0.0:" + port)
}
