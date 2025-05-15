package handlers

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// APIKeyMiddleware creates a middleware for API key authentication
func APIKeyMiddleware(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip for UI pages
		if strings.HasPrefix(c.Request.URL.Path, "/ui/") {
			c.Next()
			return
		}

		// Get API key from header
		key := c.GetHeader("X-API-Key")
		if key == "" {
			// Also check query parameter for convenience
			key = c.Query("api_key")
		}

		// Verify API key
		if key != apiKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// UIAuthMiddleware creates a middleware for UI authentication
// In a real production system, you'd want a more robust auth system
func UIAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip auth for login page
		if c.Request.URL.Path == "/ui/login" {
			c.Next()
			return
		}

		// Check for API key in cookie
		apiKey, err := c.Cookie("api_key")
		if err != nil || apiKey == "" {
			// Redirect to login page
			c.Redirect(http.StatusFound, "/ui/login")
			c.Abort()
			return
		}

		// Get API key from config for comparison
		expectedAPIKey := os.Getenv("API_KEY")
		if expectedAPIKey == "" {
			expectedAPIKey = "changeme" // Default from .env
		}

		// Verify API key
		if apiKey != expectedAPIKey {
			// Invalid API key, clear cookie and redirect to login
			c.SetCookie("api_key", "", -1, "/", "", false, true)
			c.Redirect(http.StatusFound, "/ui/login")
			c.Abort()
			return
		}

		c.Next()
	}
}
