package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return dir
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Starting server on port %s", port)
	log.Printf("Working directory: %s", getCurrentDir())
	
	// Check if directories exist
	if _, err := os.Stat("templates"); os.IsNotExist(err) {
		log.Fatalf("Templates directory not found!")
	}
	if _, err := os.Stat("static"); os.IsNotExist(err) {
		log.Fatalf("Static directory not found!")
	}
	log.Println("Directories found: templates, static")
	
	// Parse template
	templatePath := filepath.Join("templates", "index.html")
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		log.Fatalf("Template file not found: %s", templatePath)
	}
	
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}
	log.Printf("Template loaded successfully from: %s", templatePath)
	
	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("Static file server configured for /static/")
	
	// Main page handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		
		if r.URL.Path != "/" {
			log.Printf("Path not found: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, nil); err != nil {
			log.Printf("Template execution error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		log.Println("Template executed successfully")
	})
	
	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ok","service":"kvadrat-website"}`)
	})
	
	log.Printf("Server starting on http://0.0.0.0:%s", port)
	log.Println("Routes configured: /, /health, /static/*")
	
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}