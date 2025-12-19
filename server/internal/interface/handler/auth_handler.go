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
	otpService *usecase.OTPService // Dependency on OTPService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(otpService *usecase.OTPService) *AuthHandler {
	return &AuthHandler{otpService: otpService}
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
