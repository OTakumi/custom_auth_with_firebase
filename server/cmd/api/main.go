package main

import (
	"context"
	"log"
	"net/http"

	"custom_auth_api/internal/infrastructure/firebase"
	"custom_auth_api/internal/infrastructure/persistence"
	"custom_auth_api/internal/interface/handler"
	"custom_auth_api/internal/usecase"

	"github.com/gin-gonic/gin"
)

const serverPort = ":8000"

func main() {
	ctx := context.Background()

	// Initialize Firebase and Firestore client
	firestoreClient, err := firebase.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	// Initialize OTPRepository
	otpRepo := persistence.NewOTPRepository(firestoreClient)

	// Initialize OTPService
	otpService := usecase.NewOTPService(otpRepo)

	// Create the auth handler, injecting the OTPService
	authHandler := handler.NewAuthHandler(otpService)

	// Create a new Gin router
	router := gin.Default()

	// Register the routes
	router.POST("/auth/otp", authHandler.RequestOTP)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	// Start the server
	log.Printf("Server starting on port %s", serverPort)

	err = router.Run(serverPort)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
