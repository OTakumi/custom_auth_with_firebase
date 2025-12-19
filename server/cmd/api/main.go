package main

import (
	"custom_auth_api/internal/interface/handler"
	"log"
	"net/http" // Add this import

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new Gin router
	router := gin.Default()

	// Create the auth handler
	authHandler := handler.NewAuthHandler()

	// Register the routes
	router.POST("/auth/otp", authHandler.RequestOTP)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{ // Changed 200 to http.StatusOK
			"message": "OK",
		})
	})

	// Start the server
	log.Println("Server starting on port 8080")

	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
