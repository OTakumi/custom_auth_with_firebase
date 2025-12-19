package main

import (
	"log"
	"net/http"

	"custom_auth_api/internal/interface/handler"
	"custom_auth_api/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new Gin router
	router := gin.Default()

	// Initialize OTPService
	otpService := usecase.NewOTPService()

	// Create the auth handler, injecting the OTPService
	authHandler := handler.NewAuthHandler(otpService)

	// Register the routes
	router.POST("/auth/otp", authHandler.RequestOTP)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
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
