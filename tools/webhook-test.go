package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"encoding/json"
	"bytes"
	"time"
)

// Simple webhook test server that logs received webhooks to console
func main() {
	port := "8081"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		// Pretty print the JSON
		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, body, "", "  ")
		if err != nil {
			log.Printf("Error parsing JSON: %v", err)
			log.Printf("Raw payload: %s", body)
		} else {
			log.Printf("Received webhook at %s:", time.Now().Format("2006-01-02 15:04:05"))
			log.Printf("%s", prettyJSON.String())
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"success": true}`)
	})

	log.Printf("Starting webhook test server on port %s...", port)
	log.Printf("Use this URL in WhatsApp Gateway webhook settings: http://localhost:%s/webhook", port)
	log.Println("Press Ctrl+C to stop")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
