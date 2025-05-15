package whatsapp

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ClientManager manages multiple WhatsApp clients
type ClientManager struct {
	clients       map[string]*Client
	defaultClient string
	dataDir       string
	mutex         sync.RWMutex
	saveTimer     *time.Timer
}

// NewClientManager creates a new client manager
func NewClientManager(dataDir string) *ClientManager {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create data directory: %v", err))
	}

	cm := &ClientManager{
		clients: make(map[string]*Client),
		dataDir: dataDir,
	}

	// Set up periodic state saving
	cm.saveTimer = time.AfterFunc(5*time.Minute, cm.periodicSave)

	return cm
}

// periodicSave saves all client states periodically
func (cm *ClientManager) periodicSave() {
	defer cm.saveTimer.Reset(5 * time.Minute)

	if err := cm.SaveClients(); err != nil {
		fmt.Printf("Warning: Failed to save clients: %v\n", err)
	}
}

// LoadClients loads saved clients from disk
func (cm *ClientManager) LoadClients() error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Read clients directory
	entries, err := os.ReadDir(cm.dataDir)
	if err != nil {
		return fmt.Errorf("failed to read data directory: %w", err)
	}

	// Load each client
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		clientID := entry.Name()
		stateFile := filepath.Join(cm.dataDir, clientID, "state.json")

		// Check if state file exists
		if _, err := os.Stat(stateFile); os.IsNotExist(err) {
			continue
		}

		// Read state file
		data, err := os.ReadFile(stateFile)
		if err != nil {
			fmt.Printf("Warning: Failed to read state file for %s: %v\n", clientID, err)
			continue
		}

		// Parse state
		var state ClientState
		if err := json.Unmarshal(data, &state); err != nil {
			fmt.Printf("Warning: Failed to parse state for %s: %v\n", clientID, err)
			continue
		}

		// Create client
		client, err := NewClient(clientID, cm.dataDir)
		if err != nil {
			fmt.Printf("Warning: Failed to create client %s: %v\n", clientID, err)
			continue
		}

		// Add to map
		cm.clients[clientID] = client

		// Check for default client flag
		defaultFile := filepath.Join(cm.dataDir, "default_client")
		if data, err := os.ReadFile(defaultFile); err == nil {
			cm.defaultClient = string(data)
		}

		// Connect if previously connected
		if state.Status == StatusConnected || state.Connected {
			go func(c *Client) {
				if err := c.Connect(); err != nil {
					fmt.Printf("Warning: Failed to connect client %s: %v\n", c.ID, err)
				}
			}(client)
		}
	}

	return nil
}

// SaveClients saves all client states to disk
func (cm *ClientManager) SaveClients() error {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	// Save each client
	for _, client := range cm.clients {
		if err := client.SaveState(); err != nil {
			fmt.Printf("Warning: Failed to save state for %s: %v\n", client.ID, err)
		}
	}

	// Save default client
	if cm.defaultClient != "" {
		defaultFile := filepath.Join(cm.dataDir, "default_client")
		if err := os.WriteFile(defaultFile, []byte(cm.defaultClient), 0644); err != nil {
			return fmt.Errorf("failed to save default client: %w", err)
		}
	}

	return nil
}

// CreateClient creates a new WhatsApp client
func (cm *ClientManager) CreateClient(id string) (*Client, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Check if ID is empty
	if id == "" {
		return nil, errors.New("client ID cannot be empty")
	}

	// Check if ID already exists
	if _, exists := cm.clients[id]; exists {
		return nil, fmt.Errorf("client with ID %s already exists", id)
	}

	// Create client
	client, err := NewClient(id, cm.dataDir)
	if err != nil {
		return nil, err
	}

	// Add to map
	cm.clients[id] = client

	// Set as default if first client
	if len(cm.clients) == 1 && cm.defaultClient == "" {
		cm.defaultClient = id
	}

	// Save state
	if err := client.SaveState(); err != nil {
		fmt.Printf("Warning: Failed to save initial state for %s: %v\n", id, err)
	}

	return client, nil
}

// GetClient gets a client by ID
func (cm *ClientManager) GetClient(id string) (*Client, error) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	// If ID is empty, use default client
	isDefault := false
	if id == "" {
		if cm.defaultClient == "" {
			return nil, errors.New("no default client set")
		}
		id = cm.defaultClient
		isDefault = true
	} else if id == cm.defaultClient {
		isDefault = true
	}

	// Check if client exists
	client, exists := cm.clients[id]
	if !exists {
		return nil, fmt.Errorf("client %s not found", id)
	}

	// Update IsDefault flag on the client state (it won't be persisted)
	if client != nil {
		client.mutex.Lock()
		defer client.mutex.Unlock()
		// Set the IsDefault flag in the next GetState call
		client.isDefault = isDefault
	}

	return client, nil
}

// DeleteClient deletes a client
func (cm *ClientManager) DeleteClient(id string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Check if client exists
	client, exists := cm.clients[id]
	if !exists {
		return fmt.Errorf("client %s not found", id)
	}

	// First, manually disconnect and logout to ensure proper cleanup
	if client.client != nil {
		// Disconnect if connected
		if client.client.IsConnected() {
			client.client.Disconnect()
		}
		
		// Try to logout if needed
		if client.client.IsLoggedIn() {
			_ = client.client.Logout()
		}
	}

	// Close client
	if err := client.Close(); err != nil {
		return err
	}

	// Remove from map
	delete(cm.clients, id)

	// If it was the default client, unset default
	if cm.defaultClient == id {
		cm.defaultClient = ""
		
		// If there are other clients, set the first one as default
		if len(cm.clients) > 0 {
			for newDefault := range cm.clients {
				cm.defaultClient = newDefault
				break
			}
		}
	}

	// Wait a moment to ensure all handles are closed
	time.Sleep(1 * time.Second)

	// Remove client directory
	clientDir := filepath.Join(cm.dataDir, id)
	if err := os.RemoveAll(clientDir); err != nil {
		// Log but don't return error - we've already removed from memory
		fmt.Printf("Warning: Failed to remove client directory for %s: %v\n", id, err)
	}

	return nil
}

// SetDefaultClient sets the default client
func (cm *ClientManager) SetDefaultClient(id string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Check if client exists
	if _, exists := cm.clients[id]; !exists {
		return fmt.Errorf("client %s not found", id)
	}

	// Set as default
	cm.defaultClient = id

	// Save to file
	defaultFile := filepath.Join(cm.dataDir, "default_client")
	if err := os.WriteFile(defaultFile, []byte(id), 0644); err != nil {
		return fmt.Errorf("failed to save default client: %w", err)
	}

	return nil
}

// GetDefaultClient gets the default client ID
func (cm *ClientManager) GetDefaultClient() string {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	return cm.defaultClient
}

// ListClients lists all clients
func (cm *ClientManager) ListClients() []ClientState {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	// Create list
	states := make([]ClientState, 0, len(cm.clients))
	for id, client := range cm.clients {
		state := client.GetState()
		// Set IsDefault flag
		state.IsDefault = (id == cm.defaultClient)
		states = append(states, state)
	}

	return states
}

// Close closes all clients
func (cm *ClientManager) Close() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Stop save timer
	if cm.saveTimer != nil {
		cm.saveTimer.Stop()
	}

	// Save all clients before closing
	for _, client := range cm.clients {
		if err := client.SaveState(); err != nil {
			fmt.Printf("Warning: Failed to save state for %s: %v\n", client.ID, err)
		}
		
		if err := client.Close(); err != nil {
			fmt.Printf("Warning: Failed to close client %s: %v\n", client.ID, err)
		}
	}

	// Clear map
	cm.clients = make(map[string]*Client)
}
