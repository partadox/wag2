package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"go-simple-whatsapp-gateway2/config"
	"go-simple-whatsapp-gateway2/handlers"
	"go-simple-whatsapp-gateway2/whatsapp"
)

func main() {
	// Load .env file if exists
	_ = godotenv.Load()

	// Parse command line flags
	configFile := flag.String("config", "", "Path to config file")
	flag.Parse()

	// Initialize configuration
	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup client manager
	clientManager := whatsapp.NewClientManager(cfg.WhatsappDataDir)
	defer clientManager.Close()

	// Load saved clients
	if err := clientManager.LoadClients(); err != nil {
		log.Printf("Warning: Failed to load saved clients: %v", err)
	}

	// Setup router
	router := gin.Default()
	
	// Load templates with relative path
	router.LoadHTMLGlob("./templates/*")
	
	// Static files with relative path
	router.Static("/static", "./static")

	// Setup handlers
	handlers.RegisterHandlers(router, clientManager, cfg.APIKey)

	// Add debug logging
	log.Printf("Config: ListenAddr=%s, API Key=%s, WhatsappDataDir=%s", cfg.ListenAddr, cfg.APIKey, cfg.WhatsappDataDir)
	log.Printf("Template path: %s", "D:/Dev/go-simple-whatsapp-gateway2/templates/*")
	log.Printf("Static file path: %s", "D:/Dev/go-simple-whatsapp-gateway2/static")

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s", cfg.ListenAddr)
		if err := router.Run(cfg.ListenAddr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	
	// Save client states before exit
	if err := clientManager.SaveClients(); err != nil {
		log.Printf("Warning: Failed to save clients: %v", err)
	}

	log.Println("Server exited")
}
