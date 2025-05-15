package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds the application configuration
type Config struct {
	ListenAddr      string `json:"listen_addr"`
	APIKey          string `json:"api_key"`
	WhatsappDataDir string `json:"whatsapp_data_dir"`
}

// Load reads configuration from a file or environment variables
func Load(configFile string) (*Config, error) {
	// Default configuration
	cfg := &Config{
		ListenAddr:      ":8080",
		APIKey:          "changeme",
		WhatsappDataDir: "./whatsapp-data",
	}

	// Load from config file if provided
	if configFile != "" {
		if err := loadFromFile(configFile, cfg); err != nil {
			return nil, err
		}
	}

	// Override with environment variables if present
	if addr := os.Getenv("LISTEN_ADDR"); addr != "" {
		cfg.ListenAddr = addr
	}
	if key := os.Getenv("API_KEY"); key != "" {
		cfg.APIKey = key
	}
	if dir := os.Getenv("WHATSAPP_DATA_DIR"); dir != "" {
		cfg.WhatsappDataDir = dir
	}

	// Ensure the WhatsApp data directory exists
	if err := os.MkdirAll(cfg.WhatsappDataDir, 0755); err != nil {
		return nil, err
	}

	return cfg, nil
}

// loadFromFile loads configuration from a JSON file
func loadFromFile(filename string, cfg *Config) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, cfg)
}

// SaveToFile saves the configuration to a JSON file
func (cfg *Config) SaveToFile(filename string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(filename, data, 0644)
}
