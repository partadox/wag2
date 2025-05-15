package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"go-simple-whatsapp-gateway2/whatsapp"
)

// UIHandler handles UI endpoints
type UIHandler struct {
	clientManager *whatsapp.ClientManager
}

// NewUIHandler creates a new UI handler
func NewUIHandler(clientManager *whatsapp.ClientManager) *UIHandler {
	return &UIHandler{
		clientManager: clientManager,
	}
}

// RegisterRoutes registers the UI routes
func (h *UIHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/", h.redirectToDashboard)
	router.GET("/dashboard", h.dashboard)
	router.GET("/clients", h.clients)
	router.GET("/clients/:id", h.clientDetail)
	router.GET("/qrcode/:id", h.qrCode)
	router.GET("/phonepairing/:id", h.phonePairing)
	router.GET("/sendmessage/:id", h.sendMessage)
	router.GET("/test", h.testPage) // Added test route
	router.GET("/login", h.loginPage)
	router.POST("/login", h.login)
	router.GET("/logout", h.logout)
}

// redirectToDashboard redirects to the dashboard
func (h *UIHandler) redirectToDashboard(c *gin.Context) {
	c.Redirect(http.StatusFound, "/ui/dashboard")
}

// dashboard renders the dashboard page
func (h *UIHandler) dashboard(c *gin.Context) {
	clients := h.clientManager.ListClients()
	defaultClient := h.clientManager.GetDefaultClient()

	c.HTML(http.StatusOK, "dashboard_alt.html", gin.H{
		"Title":         "Dashboard",
		"Clients":       clients,
		"DefaultClient": defaultClient,
	})
}

// clients renders the clients management page
func (h *UIHandler) clients(c *gin.Context) {
	clients := h.clientManager.ListClients()
	defaultClient := h.clientManager.GetDefaultClient()

	c.HTML(http.StatusOK, "clients_alt.html", gin.H{
		"Title":         "Client Management",
		"Clients":       clients,
		"DefaultClient": defaultClient,
	})
}

// clientDetail renders the client detail page
func (h *UIHandler) clientDetail(c *gin.Context) {
	id := c.Param("id")
	client, err := h.clientManager.GetClient(id)
	if err != nil {
		c.Redirect(http.StatusFound, "/ui/clients")
		return
	}

	c.HTML(http.StatusOK, "client_detail_alt.html", gin.H{
		"Title":         "Client Details",
		"Client":        client.GetState(),
		"DefaultClient": h.clientManager.GetDefaultClient(),
	})
}

// qrCode renders the QR code page
func (h *UIHandler) qrCode(c *gin.Context) {
	id := c.Param("id")
	client, err := h.clientManager.GetClient(id)
	if err != nil {
		c.Redirect(http.StatusFound, "/ui/clients")
		return
	}

	c.HTML(http.StatusOK, "qrcode_alt2.html", gin.H{
		"Title":  "QR Code Authentication",
		"Client": client.GetState(),
	})
}

// phonePairing renders the phone pairing page
func (h *UIHandler) phonePairing(c *gin.Context) {
	id := c.Param("id")
	client, err := h.clientManager.GetClient(id)
	if err != nil {
		c.Redirect(http.StatusFound, "/ui/clients")
		return
	}

	c.HTML(http.StatusOK, "phonepairing.html", gin.H{
		"Title":  "Phone Pairing",
		"Client": client.GetState(),
	})
}

// sendMessage renders the send message page
func (h *UIHandler) sendMessage(c *gin.Context) {
	id := c.Param("id")
	client, err := h.clientManager.GetClient(id)
	if err != nil {
		c.Redirect(http.StatusFound, "/ui/clients")
		return
	}

	c.HTML(http.StatusOK, "sendmessage_alt.html", gin.H{
		"Title":  "Send Message",
		"Client": client.GetState(),
	})
}

// testPage renders a test page to verify templates and assets are loading
func (h *UIHandler) testPage(c *gin.Context) {
	c.HTML(http.StatusOK, "test_alt.html", gin.H{
		"Title":     "System Test",
		"GoVersion": "Go 1.21", // You could get the actual Go version if needed
		"Timestamp": time.Now().Format("2006-01-02 15:04:05"),
	})
}

// loginPage renders the login page
func (h *UIHandler) loginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login_alt.html", gin.H{
		"Title": "Login",
	})
}

// login processes login requests
func (h *UIHandler) login(c *gin.Context) {
	// Get API key from form
	apiKey := c.PostForm("api_key")
	
	// Get remember me
	remember := c.PostForm("remember") == "1"
	
	// Get API key from config for comparison
	expectedAPIKey := os.Getenv("API_KEY")
	if expectedAPIKey == "" {
		expectedAPIKey = "changeme" // Default from .env
	}
	
	// Verify API key
	if apiKey != expectedAPIKey {
		c.HTML(http.StatusOK, "login_alt.html", gin.H{
			"Title": "Login",
			"Error": "Invalid API Key",
		})
		return
	}
	
	// Set cookie
	expiration := 3600 // 1 hour by default
	if remember {
		expiration = 3600 * 24 // 24 hours if remember me is checked
	}
	
	c.SetCookie("api_key", apiKey, expiration, "/", "", false, true)
	
	// Redirect to dashboard
	c.Redirect(http.StatusFound, "/ui/dashboard")
}

// logout handles logout by clearing cookies
func (h *UIHandler) logout(c *gin.Context) {
	// Clear cookie
	c.SetCookie("api_key", "", -1, "/", "", false, true)
	
	// Redirect to login page
	c.Redirect(http.StatusFound, "/ui/login")
}