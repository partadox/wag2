package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go-simple-whatsapp-gateway2/whatsapp"
)

// WhatsAppHandler handles legacy single-client API endpoints
type WhatsAppHandler struct {
	clientManager *whatsapp.ClientManager
}

// NewWhatsAppHandler creates a new WhatsApp handler
func NewWhatsAppHandler(clientManager *whatsapp.ClientManager) *WhatsAppHandler {
	return &WhatsAppHandler{
		clientManager: clientManager,
	}
}

// RegisterRoutes registers the legacy API routes
func (h *WhatsAppHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/status", h.getStatus)
	router.GET("/qr", h.generateQR)
	router.POST("/pair", h.pairPhone)
	router.GET("/paircode", h.getPairingCode)
	router.POST("/send", h.sendMessage)
	router.POST("/connect", h.connect)
	router.POST("/disconnect", h.disconnect)
	router.POST("/logout", h.logout)
}

// getStatus gets the status of the default client
func (h *WhatsAppHandler) getStatus(c *gin.Context) {
	client, err := h.clientManager.GetClient("")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, client.GetState())
}

// generateQR generates a QR code for the default client
func (h *WhatsAppHandler) generateQR(c *gin.Context) {
	client, err := h.clientManager.GetClient("")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	qrCode, err := client.GenerateQR()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"qr_code": qrCode})
}

// pairPhone pairs the default client with a phone number (currently not supported)
func (h *WhatsAppHandler) pairPhone(c *gin.Context) {
	client, err := h.clientManager.GetClient("")
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

// getPairingCode gets the pairing code for the default client (currently not supported)
func (h *WhatsAppHandler) getPairingCode(c *gin.Context) {
	client, err := h.clientManager.GetClient("")
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

// sendMessage sends a message from the default client
func (h *WhatsAppHandler) sendMessage(c *gin.Context) {
	client, err := h.clientManager.GetClient("")
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

// connect connects the default client
func (h *WhatsAppHandler) connect(c *gin.Context) {
	client, err := h.clientManager.GetClient("")
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

// disconnect disconnects the default client
func (h *WhatsAppHandler) disconnect(c *gin.Context) {
	client, err := h.clientManager.GetClient("")
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

// logout logs out the default client
func (h *WhatsAppHandler) logout(c *gin.Context) {
	client, err := h.clientManager.GetClient("")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := client.Logout(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, client.GetState())
}
