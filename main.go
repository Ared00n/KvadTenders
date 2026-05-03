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
	
	// Parse template
	templatePath := filepath.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}
	log.Println("Template loaded successfully")
	
	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	
	// Main page handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, nil); err != nil {
			log.Printf("Template execution error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
	
	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ok","service":"kvadrat-website"}`)
	})
	
	log.Printf("Server starting on port %s", port)
	
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}