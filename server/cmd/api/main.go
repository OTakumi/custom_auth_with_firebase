package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"custom_auth_api/internal/infrastructure/emailsender"
	"custom_auth_api/internal/infrastructure/firebase"
	"custom_auth_api/internal/infrastructure/persistence"
	"custom_auth_api/internal/interface/handler"
	"custom_auth_api/internal/interface/middleware"
	"custom_auth_api/internal/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

const serverPort = ":8000"

func main() {
	ctx := context.Background()

	// Initialize Firebase and Firestore client
	firestoreClient, authClient, err := firebase.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}
	// Ensure Firestore client is properly closed on shutdown
	defer func() {
		err := firestoreClient.Close()
		if err != nil {
			log.Printf("Error closing Firestore client: %v", err)
		}
	}()

	// Use authClient for other services
	authService := usecase.NewAuthService(authClient)

	// Initialize dependencies
	otpRepo := persistence.NewOTPRepository(firestoreClient)
	emailSender := emailsender.NewDummyEmailSender()

	// Initialize OTPService
	otpService := usecase.NewOTPService(otpRepo, emailSender)

	// Create the handlers, injecting the OTPService and AuthService
	otpRequestHandler := handler.NewOTPRequestHandler(otpService, authService)
	otpVerifyHandler := handler.NewOTPVerifyHandler(otpService, authService)

	// Create a new Gin router
	router := gin.Default()

	// CORS configuration
	corsConfig := cors.DefaultConfig()

	env := os.Getenv("ENV")
	if env == "production" {
		// Production: Allow only specified origins
		allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
		if allowedOrigins != "" {
			corsConfig.AllowOrigins = strings.Split(allowedOrigins, ",")
			log.Printf("CORS: Allowing origins: %v", corsConfig.AllowOrigins)
		} else {
			log.Fatal("ALLOWED_ORIGINS environment variable is required in production")
		}
	} else {
		// Development: Allow all origins
		corsConfig.AllowAllOrigins = true

		log.Println("CORS: Allowing all origins (development mode)")
	}

	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Content-Type", "Authorization"}
	router.Use(cors.New(corsConfig))

	// Rate limiting configuration
	// 5 requests per minute per IP for authentication endpoints
	rateLimiter := middleware.NewIPRateLimiter(rate.Every(time.Minute/5), 5)
	// Start cleanup routine to prevent memory leak
	go rateLimiter.CleanupExpiredLimiters(10 * time.Minute)

	// Register the routes with rate limiting
	authGroup := router.Group("/auth")
	authGroup.Use(middleware.RateLimitMiddleware(rateLimiter))
	{
		authGroup.POST("/otp", otpRequestHandler.RequestOTP)
		authGroup.POST("/verify", otpVerifyHandler.VerifyOTP)
	}

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
