package main

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"koalatranslate/internal/auth"
	"koalatranslate/internal/db"
	"koalatranslate/internal/httpapi"
	"koalatranslate/internal/projects"
)

//go:embed all:build
var embeddedFiles embed.FS

func main() {
	log.Println("Starting KoalaTranslate server...")

	// 1. Read configurations from environment variables
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "koalatranslate.db"
	}

	adminUser := os.Getenv("ADMIN_USERNAME")
	if adminUser == "" {
		adminUser = "admin"
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "adminpassword"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 2. Initialize Database and run migrations
	database, err := db.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer database.Close()

	// 3. Seed initial admin user and default project
	if err := auth.SeedAdmin(database, adminUser, adminPassword); err != nil {
		log.Fatalf("Failed to seed admin user: %v", err)
	}
	log.Printf("Admin user '%s' initialized.\n", adminUser)

	if err := projects.SeedProjects(database); err != nil {
		log.Fatalf("Failed to seed projects: %v", err)
	}
	log.Println("Default projects seeded.")

	// 4. Set up HTTP router
	mux := http.NewServeMux()

	// Register API endpoints
	apiServer := httpapi.NewServer(database)
	apiServer.RegisterRoutes(mux)

	// Serve SvelteKit static files (SPA fallback handler)
	subFS, err := fs.Sub(embeddedFiles, "build")
	if err != nil {
		log.Fatalf("Failed to get sub FS of build: %v", err)
	}

	// Check if we are running in Dev CORS/Proxy mode
	corsEnabled := os.Getenv("DEV_CORS") == "true"

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// API calls are already matched and served by registerRoutes,
		// but check /api/ prefix just in case
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		// Enable CORS in development mode
		if corsEnabled {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Cookie")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		// Try opening the file in static FS
		f, err := subFS.Open(path)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				// Fallback to index.html for client-side routing
				indexFile, err := subFS.Open("index.html")
				if err != nil {
					w.Header().Set("Content-Type", "text/plain; charset=utf-8")
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprint(w, "Index file not found. Make sure the SvelteKit build is completed.")
					return
				}
				defer indexFile.Close()
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				_, _ = io.Copy(w, indexFile)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		// Retrieve content type
		stat, err := f.Stat()
		if err == nil {
			contentType := mimeTypeByExtension(filepath.Ext(path))
			w.Header().Set("Content-Type", contentType)
			w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
		}

		_, _ = io.Copy(w, f)
	})

	// Wrap serve multiplexer in simple logging middleware
	loggingMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Never log ADMIN_PASSWORD, SESSION_SECRET, GITHUB_TOKEN, cookies, or session IDs
		urlPath := r.URL.Path
		method := r.Method
		log.Printf("%s %s\n", method, urlPath)
		mux.ServeHTTP(w, r)
	})

	log.Printf("Listening on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, loggingMux); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}

// Simple lookup helper for common static file types
func mimeTypeByExtension(ext string) string {
	switch strings.ToLower(ext) {
	case ".html":
		return "text/html; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".js":
		return "application/javascript; charset=utf-8"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	case ".json":
		return "application/json"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	default:
		return "application/octet-stream"
	}
}

// Database helper
func getDBConnection(db *sql.DB) *sql.DB {
	return db
}
