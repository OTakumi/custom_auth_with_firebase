package main

import (
	"context"
	"log"

	"custom_auth_api/internal/config"
	"custom_auth_api/internal/infrastructure/emailsender"
	"custom_auth_api/internal/infrastructure/firebase"
	"custom_auth_api/internal/infrastructure/persistence"
	"custom_auth_api/internal/interface/handler"
	"custom_auth_api/internal/interface/router"
	"custom_auth_api/internal/usecase"
)

func main() {
	ctx := context.Background()

	// Load environment configuration
	env, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("Failed to load environment configuration: %v", err)
	}

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

	// Initialize services
	authService := usecase.NewAuthService(authClient)
	otpSessionRepo := persistence.NewOTPSessionRepository(firestoreClient)
	emailSender := emailsender.NewDummyEmailSender()
	otpService := usecase.NewOTPService(otpSessionRepo, emailSender)

	// Initialize handlers
	handlers := &router.Handlers{
		OTPRequest: handler.NewOTPRequestHandler(otpService, authService),
		OTPVerify:  handler.NewOTPVerifyHandler(otpService, authService),
	}

	// Setup router with all middleware and routes
	r := router.NewRouter(env, handlers)

	// Start the server
	serverAddr := ":" + env.Port
	log.Printf("Server starting on port %s (environment: %s)", serverAddr, env.Environment)

	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Server failed to start: %v", err) //nolint:gocritic // log.Fatalf is intentional
	}
}
