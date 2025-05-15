package main

import (
	"fmt"
	"os"
)

func main() {
	// Delete test-webhook-server.go
	if err := os.Remove("test-webhook-server.go"); err != nil {
		fmt.Printf("Failed to delete test-webhook-server.go: %v\n", err)
	} else {
		fmt.Println("Successfully deleted test-webhook-server.go")
	}

	// Delete webhook-test.go
	if err := os.Remove("webhook-test.go"); err != nil {
		fmt.Printf("Failed to delete webhook-test.go: %v\n", err)
	} else {
		fmt.Println("Successfully deleted webhook-test.go")
	}
}
