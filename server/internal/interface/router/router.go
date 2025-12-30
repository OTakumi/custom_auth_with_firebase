package router

import (
	"log"
	"net/http"
	"time"

	"custom_auth_api/internal/config"
	"custom_auth_api/internal/interface/handler"
	"custom_auth_api/internal/interface/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Handlers holds all HTTP handlers for dependency injection.
type Handlers struct {
	OTPRequest *handler.OTPRequestHandler
	OTPVerify  *handler.OTPVerifyHandler
}

// NewRouter creates and configures a new Gin router with all middleware and routes.
// This function encapsulates all router setup logic including CORS, rate limiting, and route registration.
func NewRouter(env *config.Env, handlers *Handlers) *gin.Engine {
	router := gin.Default()

	// Setup CORS middleware
	router.Use(setupCORS(env))

	// Setup rate limiting
	rateLimiter := setupRateLimiter(env)

	// Register routes
	registerRoutes(router, rateLimiter, handlers)

	return router
}

// setupCORS configures CORS middleware based on environment settings.
// In production: Only allows specified origins from ALLOWED_ORIGINS env var
// In development: Allows all origins
func setupCORS(env *config.Env) gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()

	if env.IsProduction() {
		// Production: Allow only specified origins
		corsConfig.AllowOrigins = env.AllowedOrigins
		log.Printf("CORS: Allowing origins: %v", corsConfig.AllowOrigins)
	} else {
		// Development: Allow all origins
		corsConfig.AllowAllOrigins = true
		log.Println("CORS: Allowing all origins (development mode)")
	}

	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Content-Type", "Authorization"}

	return cors.New(corsConfig)
}

// setupRateLimiter creates and configures the IP-based rate limiter.
// Starts a background cleanup routine to prevent memory leaks.
func setupRateLimiter(env *config.Env) *middleware.IPRateLimiter {
	requestsPerMinute := env.RateLimitRequestsPerMinute
	cleanupInterval := time.Duration(env.RateLimitCleanupIntervalMinutes) * time.Minute

	rateLimiter := middleware.NewIPRateLimiter(
		rate.Every(time.Minute/time.Duration(requestsPerMinute)),
		requestsPerMinute,
	)

	// Start cleanup routine to prevent memory leak
	go rateLimiter.CleanupExpiredLimiters(cleanupInterval)

	return rateLimiter
}

// registerRoutes registers all application routes with appropriate middleware.
func registerRoutes(router *gin.Engine, rateLimiter *middleware.IPRateLimiter, handlers *Handlers) {
	// Health check endpoint (no rate limiting)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	// Authentication endpoints with rate limiting
	authGroup := router.Group("/auth")
	authGroup.Use(middleware.RateLimitMiddleware(rateLimiter))
	{
		authGroup.POST("/otp", handlers.OTPRequest.RequestOTP)
		authGroup.POST("/verify", handlers.OTPVerify.VerifyOTP)
	}
}
