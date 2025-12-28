package handler

import (
	"log"
	"net/http"

	"custom_auth_api/internal/domain/vo/email"
	"custom_auth_api/internal/usecase"

	"github.com/gin-gonic/gin"
)

// OTPRequestHandler handles OTP request related operations.
//
// Responsibilities:
// - Handle POST /auth/otp endpoint
// - Validate email format
// - Check user existence before generating OTP
// - Generate and send OTP to registered users.
type OTPRequestHandler struct {
	otpService  *usecase.OTPService
	authService *usecase.AuthService
}

// NewOTPRequestHandler creates a new OTPRequestHandler.
func NewOTPRequestHandler(otpService *usecase.OTPService, authService *usecase.AuthService) *OTPRequestHandler {
	return &OTPRequestHandler{
		otpService:  otpService,
		authService: authService,
	}
}

// RequestOTP is a handler for generating an OTP.
func (h *OTPRequestHandler) RequestOTP(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})

		return
	}

	// Validate email format using the value object
	_, err = email.NewEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	// Check if user exists in Firebase Auth before generating OTP
	_, err = h.authService.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		// Use generic error message to prevent email enumeration attacks
		log.Printf("Authentication failed for OTP request: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})

		return
	}

	// Generate and save OTP using the service
	_, err = h.otpService.GenerateAndSendOTP(c.Request.Context(), req.Email)
	if err != nil {
		log.Printf("Error generating and saving OTP for %s: %v", req.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate and save OTP"})

		return
	}

	// Return success message without exposing OTP
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully."})
}
