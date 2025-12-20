package handler

import (
	"log"
	"net/http"

	"custom_auth_api/internal/domain/vo/email"
	"custom_auth_api/internal/usecase"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication related requests.
type AuthHandler struct {
	otpService  *usecase.OTPService  // Dependency on OTPService
	authService *usecase.AuthService // Dependency on AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(otpService *usecase.OTPService, authService *usecase.AuthService) *AuthHandler {
	return &AuthHandler{otpService: otpService, authService: authService}
}

// RequestOTP is a handler for generating an OTP.
func (h *AuthHandler) RequestOTP(c *gin.Context) {
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

	// Generate and save OTP using the service
	otp, err := h.otpService.GenerateAndSendOTP(c.Request.Context(), req.Email)
	if err != nil {
		log.Printf("Error generating and saving OTP for %s: %v", req.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate and save OTP"})

		return
	}

	// For now, just return a success message. The OTP is logged by the service.
	c.JSON(http.StatusOK, gin.H{"message": "OTP generation request received", "otp": otp})
}

// VerifyOTP is a handler for verifying an OTP and generating a custom token.
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
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

	// If OTP is valid, generate a custom Firebase token
	// For now, use the email as UID. In a real application, you might map this to a proper user ID.
	customToken, err := h.authService.GenerateCustomToken(c.Request.Context(), req.Email)
	if err != nil {
		log.Printf("Error generating custom token for %s: %v", req.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication token"})

		return
	}

	c.JSON(http.StatusOK, gin.H{"token": customToken})
}
