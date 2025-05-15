package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go-simple-whatsapp-gateway2/whatsapp"
)

// ClientRequest represents a client creation/update request
type ClientRequest struct {
	ID string `json:"id" binding:"required"`
}

// DefaultClientRequest represents a request to set the default client
type DefaultClientRequest struct {
	ID string `json:"id" binding:"required"`
}

// PairingRequest represents a phone pairing request
type PairingRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

// MessageRequest represents a message sending request
type MessageRequest struct {
	Recipient string `json:"recipient" binding:"required"`
	Message   string `json:"message" binding:"required"`
}

// WebhookRequest represents a webhook configuration request
type WebhookRequest struct {
	Enabled bool   `json:"enabled"`
	URL     string `json:"url"`
}

// ClientsHandler handles multi-client API endpoints
type ClientsHandler struct {
	clientManager *whatsapp.ClientManager
}

// NewClientsHandler creates a new clients handler
func NewClientsHandler(clientManager *whatsapp.ClientManager) *ClientsHandler {
	return &ClientsHandler{
		clientManager: clientManager,
	}
}

// RegisterRoutes registers the client API routes
func (h *ClientsHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/clients", h.listClients)
	router.POST("/clients", h.createClient)
	router.POST("/clients/default", h.setDefaultClient)
	router.GET("/clients/:id", h.getClient)
	router.DELETE("/clients/:id", h.deleteClient)
	router.GET("/clients/:id/qr", h.generateQR)
	router.POST("/clients/:id/pair", h.pairPhone)
	router.GET("/clients/:id/paircode", h.getPairingCode)
	router.POST("/clients/:id/send", h.sendMessage)
	router.POST("/clients/:id/connect", h.connectClient)
	router.POST("/clients/:id/disconnect", h.disconnectClient)
	router.POST("/clients/:id/logout", h.logoutClient)
	// Webhook routes
	router.GET("/clients/:id/webhook", h.getWebhookConfig)
	router.POST("/clients/:id/webhook", h.setWebhookConfig)
}

// listClients lists all clients
func (h *ClientsHandler) listClients(c *gin.Context) {
	clients := h.clientManager.ListClients()
	c.JSON(http.StatusOK, gin.H{
		"clients":        clients,
		"default_client": h.clientManager.GetDefaultClient(),
	})
}

// createClient creates a new client
func (h *ClientsHandler) createClient(c *gin.Context) {
	var req ClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	client, err := h.clientManager.CreateClient(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, client.GetState())
}

// getClient gets a client by ID
func (h *ClientsHandler) getClient(c *gin.Context) {
	id := c.Param("id")
	client, err := h.clientManager.GetClient(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, client.GetState())
}

// deleteClient deletes a client
func (h *ClientsHandler) deleteClient(c *gin.Context) {
	id := c.Param("id")
	if err := h.clientManager.DeleteClient(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// setDefaultClient sets the default client
func (h *ClientsHandler) setDefaultClient(c *gin.Context) {
	var req DefaultClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.clientManager.SetDefaultClient(req.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// generateQR generates a QR code for a client
func (h *ClientsHandler) generateQR(c *gin.Context) {
	id := c.Param("id")
	log.Printf("Generating QR code for client: %s", id) // Add logging
	
	client, err := h.clientManager.GetClient(id)
	if err != nil {
		log.Printf("Error getting client %s: %v", id, err) // Add logging
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// If the client is already logged in, return an appropriate error
	state := client.GetState()
	if state.LoggedIn {
		log.Printf("Client %s is already logged in, no need for QR code", id)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Client is already logged in. Logout first if you want to reconnect.",
			"logged_in": true,
		})
		return
	}

	// Try to generate the QR code
	qrCode, err := client.GenerateQR()
	if err != nil {
		log.Printf("Error generating QR code for client %s: %v", id, err) // Add logging
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("QR code generated successfully for client %s", id) // Add logging
	c.JSON(http.StatusOK, gin.H{"qr_code": qrCode})
}

// pairPhone pairs a client with a phone number (currently not supported)
func (h *ClientsHandler) pairPhone(c *gin.Context) {
	id := c.Param("id")
	client, err := h.clientManager.GetClient(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var req PairingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := client.PairPhone(req.PhoneNumber); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// getPairingCode gets the pairing code after a PairPhone request (currently not supported)
func (h *ClientsHandler) getPairingCode(c *gin.Context) {
	id := c.Param("id")
	client, err := h.clientManager.GetClient(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	code, err := client.GetPairingCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": code})
}

// sendMessage sends a message from a client
func (h *ClientsHandler) sendMessage(c *gin.Context) {
	id := c.Param("id")
	client, err := h.clientManager.GetClient(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var req MessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := client.SendMessage(req.Recipient, req.Message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"sent_at": time.Now(),
	})
}

// connectClient connects a client
func (h *ClientsHandler) connectClient(c *gin.Context) {
	id := c.Param("id")
	client, err := h.clientManager.GetClient(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := client.Connect(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, client.GetState())
}

// disconnectClient disconnects a client
func (h *ClientsHandler) disconnectClient(c *gin.Context) {
	id := c.Param("id")
	client, err := h.clientManager.GetClient(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := client.Disconnect(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, client.GetState())
}

// logoutClient logs out a client
func (h *ClientsHandler) logoutClient(c *gin.Context) {
	id := c.Param("id")
	log.Printf("Logging out client: %s", id)
	
	client, err := h.clientManager.GetClient(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// First disconnect, then logout
	// Disconnect silently - don't throw error if already disconnected
	client.Disconnect()

	// Now attempt formal logout
	err = client.Logout()
	if err != nil {
		log.Printf("Error logging out client %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Client %s logged out successfully", id)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// getWebhookConfig gets webhook configuration for a client
func (h *ClientsHandler) getWebhookConfig(c *gin.Context) {
	id := c.Param("id")
	client, err := h.clientManager.GetClient(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Get webhook configuration
	webhookConfig, err := client.GetWebhookConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, webhookConfig)
}

// setWebhookConfig sets webhook configuration for a client
func (h *ClientsHandler) setWebhookConfig(c *gin.Context) {
	id := c.Param("id")
	client, err := h.clientManager.GetClient(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var req WebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate URL if webhook is enabled
	if req.Enabled && req.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Webhook URL is required when enabled"})
		return
	}

	// Set webhook configuration
	err = client.SaveWebhookConfig(whatsapp.WebhookConfig{
		Enabled: req.Enabled,
		URL:     req.URL,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save client state to persist webhook configuration
	client.SaveState()

	c.JSON(http.StatusOK, gin.H{"success": true})
}
