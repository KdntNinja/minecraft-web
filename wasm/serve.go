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
		// Enhanced log: timestamp, method, path, status, duration, client, user-agent, referer
		reqTime := time.Now().Format("2006-01-02 15:04:05")
		log.Printf("\033[1;30m[%s]\033[0m \033[1;34m%s\033[0m \033[1;36m%-20s\033[0m from \033[1;33m%-15s\033[0m\n    \033[2mUser-Agent:\033[0m %s\n    \033[2mReferer:\033[0m %s",
			reqTime, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), r.Referer())

		// Security headers (improved)
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=(), fullscreen=(self)")
		// CSP: allow WASM, block inline scripts, allow 'unsafe-eval' for Go WASM, restrict everything else
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-eval'; style-src 'self'; img-src 'self' data:; connect-src 'self'; font-src 'self'; object-src 'none'; base-uri 'self'; form-action 'self'; frame-ancestors 'self';")

		// CORS headers (locked down: only allow GET, restrict origin, allow localhost and 127.0.0.1)
		allowedOrigins := []string{
			"http://localhost:" + port,
			"http://127.0.0.1:" + port,
			"https://webcraft.kdnsite.site",
		}
		origin := r.Header.Get("Origin")
		for _, ao := range allowedOrigins {
			if origin == ao {
				w.Header().Set("Access-Control-Allow-Origin", ao)
				break
			}
		}
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// MIME types (improved, add more types and fallback)
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
		case ".jpg", ".jpeg":
			w.Header().Set("Content-Type", "image/jpeg")
		case ".gif":
			w.Header().Set("Content-Type", "image/gif")
		case ".ico":
			w.Header().Set("Content-Type", "image/x-icon")
		case ".svg":
			w.Header().Set("Content-Type", "image/svg+xml")
		case ".json":
			w.Header().Set("Content-Type", "application/json")
		case ".mp3":
			w.Header().Set("Content-Type", "audio/mpeg")
		case ".wav":
			w.Header().Set("Content-Type", "audio/wav")
		case ".mp4":
			w.Header().Set("Content-Type", "video/mp4")
		case ".webm":
			w.Header().Set("Content-Type", "video/webm")
		case ".txt":
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		default:
			w.Header().Set("Content-Type", "application/octet-stream")
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

		// Directory traversal and listing prevention
		cleanPath := filepath.Clean(r.URL.Path)
		if cleanPath == "." || cleanPath == "/" {
			cleanPath = "/index.html"
		}
		if cleanPath[0] == '/' {
			cleanPath = "." + cleanPath
		}
		absPath, err := filepath.Abs(cleanPath)
		if err != nil || len(absPath) < len(currentDir) || absPath[:len(currentDir)] != currentDir {
			http.Error(w, "403 Forbidden: Invalid path", http.StatusForbidden)
			return
		}
		info, err := os.Stat(absPath)
		if err == nil && info.IsDir() {
			indexPath := filepath.Join(absPath, "index.html")
			if _, err := os.Stat(indexPath); err == nil {
				http.ServeFile(w, r, indexPath)
				return
			}
			http.Error(w, "403 Forbidden: Directory listing is disabled", http.StatusForbidden)
			return
		}
		if err != nil || info == nil {
			http.NotFound(w, r)
			return
		}

		// Serve file
		fs.ServeHTTP(w, r)

		// Log response time and status
		duration := time.Since(start)
		statusColor := "\033[1;32m[OK]\033[0m"
		if r.Method == "OPTIONS" {
			statusColor = "\033[1;36m[OPT]\033[0m"
		} else if r.URL.Path == "/favicon.ico" {
			statusColor = "\033[1;31m[404]\033[0m"
		}
		log.Printf("    %s \033[1;36m%-20s\033[0m in \033[1;35m%v\033[0m", statusColor, r.URL.Path, duration)
	})

	fmt.Printf("\n\033[1;33mStarting Webcraft server on port %s\033[0m\n", port)
	fmt.Printf("Open \033[4;36mhttp://localhost:%s\033[0m in your browser\n\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
