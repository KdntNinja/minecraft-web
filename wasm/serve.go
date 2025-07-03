package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}
	fmt.Printf("\n\033[1;32mWebcraft Static Server\033[0m\nServing files from: \033[1;36m%s\033[0m\n", currentDir)

	fs := http.FileServer(http.Dir("."))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Log the request with more detail
		log.Printf("\033[1;34m[REQ]\033[0m %s \033[1;36m%s\033[0m from \033[1;33m%s\033[0m\n        \033[2mUser-Agent:\033[0m %s\n        \033[2mReferer:\033[0m %s",
			r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), r.Referer())

		// Security headers
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")

		// CORS headers for development
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		// MIME types
		ext := filepath.Ext(r.URL.Path)
		switch ext {
		case ".wasm":
			w.Header().Set("Content-Type", "application/wasm")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		case ".html":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		case ".png":
			w.Header().Set("Content-Type", "image/png")
		case ".ico":
			w.Header().Set("Content-Type", "image/x-icon")
		case ".svg":
			w.Header().Set("Content-Type", "image/svg+xml")
		case ".json":
			w.Header().Set("Content-Type", "application/json")
		}

		// Preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Favicon
		if r.URL.Path == "/favicon.ico" {
			w.Header().Set("Content-Type", "image/x-icon")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Directory listing prevention
		filePath := "." + r.URL.Path
		info, err := os.Stat(filePath)
		if err == nil && info.IsDir() {
			indexPath := filepath.Join(filePath, "index.html")
			if _, err := os.Stat(indexPath); err == nil {
				http.ServeFile(w, r, indexPath)
				return
			}
			http.Error(w, "403 Forbidden: Directory listing is disabled", http.StatusForbidden)
			return
		}

		// Serve file
		fs.ServeHTTP(w, r)

		// Log response time
		duration := time.Since(start)
		log.Printf("\033[1;32m[OK]\033[0m Served \033[1;36m%s\033[0m in \033[1;35m%v\033[0m", r.URL.Path, duration)
	})

	fmt.Printf("\n\033[1;33mStarting Webcraft server on port %s\033[0m\n", port)
	fmt.Printf("Open \033[4;36mhttp://localhost:%s\033[0m in your browser\n\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
