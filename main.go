// Package main implements a web server for Kvadrat company website.
// Kvadrat specializes in goods supply through electronic auctions.
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Config holds application configuration
type Config struct {
	Port         string
	StaticDir    string
	TemplatesDir string
}

// Server represents the HTTP server with its dependencies
type Server struct {
	config   *Config
	template *template.Template
}

// NewServer creates a new server instance with the given configuration
func NewServer(config *Config) (*Server, error) {
	// Parse templates
	templatePath := filepath.Join(config.TemplatesDir, "index.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return &Server{
		config:   config,
		template: tmpl,
	}, nil
}

// setupRoutes configures HTTP routes for the server
func (s *Server) setupRoutes() {
	// Serve static files (CSS, images, etc.)
	staticHandler := http.StripPrefix("/static/", 
		http.FileServer(http.Dir(s.config.StaticDir)))
	http.Handle("/static/", staticHandler)

	// Main page handler
	http.HandleFunc("/", s.homeHandler)
	
	// Health check endpoint
	http.HandleFunc("/health", s.healthHandler)
}

// homeHandler handles requests to the main page
func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	// Only serve GET requests to root path
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Set content type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	
	// Execute template without data
	if err := s.template.Execute(w, nil); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// healthHandler provides a simple health check endpoint
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"ok","service":"kvadrat-website"}`)
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.setupRoutes()
	
	log.Printf("Starting server on port %s", s.config.Port)
	log.Printf("Server URL: http://localhost%s", s.config.Port)
	
	return http.ListenAndServe(s.config.Port, nil)
}

// getConfig returns application configuration from environment variables or defaults
func getConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	// Render передает порт без двоеточия, добавляем его
	if port[0] != ':' {
		port = ":" + port
	}

	return &Config{
		Port:         port,
		StaticDir:    "static",
		TemplatesDir: "templates",
	}
}

// validateDirectories checks if required directories exist
func validateDirectories(config *Config) error {
	dirs := []string{config.StaticDir, config.TemplatesDir}
	
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return fmt.Errorf("required directory does not exist: %s", dir)
		}
	}
	
	return nil
}

func main() {
	// Get configuration
	config := getConfig()
	
	// Validate required directories
	if err := validateDirectories(config); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}
	
	// Create server instance
	server, err := NewServer(config)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	
	// Start server
	log.Fatal(server.Start())
}