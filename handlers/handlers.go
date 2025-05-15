package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go-simple-whatsapp-gateway2/whatsapp"
)

// RegisterHandlers registers all the handlers
func RegisterHandlers(router *gin.Engine, clientManager *whatsapp.ClientManager, apiKey string) {
	// Middleware for API authentication
	apiAuthMiddleware := APIKeyMiddleware(apiKey)
	uiAuthMiddleware := UIAuthMiddleware()

	// API routes
	apiGroup := router.Group("/api")
	apiGroup.Use(apiAuthMiddleware)

	// Legacy single-client API
	whatsAppHandler := NewWhatsAppHandler(clientManager)
	whatsAppHandler.RegisterRoutes(apiGroup)

	// Multi-client API
	clientsHandler := NewClientsHandler(clientManager)
	clientsHandler.RegisterRoutes(apiGroup)

	// UI routes
	uiGroup := router.Group("/ui")
	uiGroup.Use(uiAuthMiddleware)

	uiHandler := NewUIHandler(clientManager)
	uiHandler.RegisterRoutes(uiGroup)

	// Redirect root to UI
	router.GET("/", func(c *gin.Context) {
		// Check if user is authenticated
		_, err := c.Cookie("api_key")
		if err != nil {
			// Not authenticated, redirect to login
			c.Redirect(http.StatusFound, "/ui/login")
			return
		}
		// Authenticated, redirect to dashboard
		c.Redirect(http.StatusFound, "/ui/dashboard")
	})
	
	// Add a test route at root level for troubleshooting
	router.GET("/test", func(c *gin.Context) {
		c.HTML(http.StatusOK, "test_alt.html", gin.H{
			"Title":     "System Test",
			"GoVersion": "Go 1.21",
			"Timestamp": time.Now().Format("2006-01-02 15:04:05"),
		})
	})
}
