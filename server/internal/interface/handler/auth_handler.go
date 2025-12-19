package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication related requests.
type AuthHandler struct{}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
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

	// For now, just log the email and return a success message.
	log.Printf("Received email: %s", req.Email)

	c.JSON(http.StatusOK, gin.H{"message": "OTP generation request received"})
}
