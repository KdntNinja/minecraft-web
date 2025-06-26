package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Get port from environment or default to 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Create file server
	fs := http.FileServer(http.Dir("."))

	// Add CORS and proper MIME type headers for WASM
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		// Set proper MIME type for WASM files
		if strings.HasSuffix(r.URL.Path, ".wasm") {
			w.Header().Set("Content-Type", "application/wasm")
		}

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			return
		}

		fs.ServeHTTP(w, r)
	})

	log.Printf("Server starting on :%s", port)
	log.Printf("Open http://localhost:%s in your browser", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
