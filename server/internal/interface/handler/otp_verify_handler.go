package handler

import (
	"log"
	"net/http"

	"custom_auth_api/internal/domain/vo/email"
	"custom_auth_api/internal/usecase"

	"github.com/gin-gonic/gin"
)

// OTPVerifyHandler handles OTP verification and token generation.
//
// Responsibilities:
// - Handle POST /auth/verify endpoint
// - Validate email format
// - Verify OTP against stored value
// - Generate Firebase custom token for authenticated users.
type OTPVerifyHandler struct {
	otpService  *usecase.OTPService
	authService *usecase.AuthService
}

// NewOTPVerifyHandler creates a new OTPVerifyHandler.
func NewOTPVerifyHandler(otpService *usecase.OTPService, authService *usecase.AuthService) *OTPVerifyHandler {
	return &OTPVerifyHandler{
		otpService:  otpService,
		authService: authService,
	}
}

// VerifyOTP is a handler for verifying an OTP and generating a custom token.
func (h *OTPVerifyHandler) VerifyOTP(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
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

	// Verify the OTP
	isValid, err := h.otpService.VerifyOTP(c.Request.Context(), req.Email, req.OTP)
	if err != nil || !isValid {
		// Log the error for internal tracking, but return a generic invalid OTP message to the client
		log.Printf("OTP verification failed for %s: %v", req.Email, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})

		return
	}

	// Check if user exists in Firebase Auth
	user, err := h.authService.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		// Use generic error message to prevent email enumeration attacks
		log.Printf("Authentication failed for OTP verification: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})

		return
	}

	// If OTP is valid and user exists, generate a custom Firebase token
	customToken, err := h.authService.GenerateCustomToken(c.Request.Context(), user.UID)
	if err != nil {
		log.Printf("Error generating custom token for %s: %v", req.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication token"})

		return
	}

	c.JSON(http.StatusOK, gin.H{"token": customToken})
}
