package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/mumugogoing/meme_bot/internal/config"
	"github.com/mumugogoing/meme_bot/pkg/meme"
	"github.com/rs/cors"
)

type MemeRequest struct {
	Template   string `json:"template"`
	ImageURL   string `json:"image_url"`
	TopText    string `json:"top_text"`
	BottomText string `json:"bottom_text"`
}

type TemplatesResponse struct {
	Templates []string `json:"templates"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Server struct {
	generator *meme.Generator
	config    *config.Config
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	generator, err := meme.NewGenerator(cfg.TemplatesDir, cfg.OutputDir)
	if err != nil {
		log.Fatalf("Failed to create meme generator: %v", err)
	}

	server := &Server{
		generator: generator,
		config:    cfg,
	}

	router := mux.NewRouter()
	
	// API routes
	router.HandleFunc("/api/templates", server.handleListTemplates).Methods("GET")
	router.HandleFunc("/api/meme", server.handleCreateMeme).Methods("POST")
	router.HandleFunc("/api/health", server.handleHealth).Methods("GET")

	// Serve static files from frontend build
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend/dist")))

	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	addr := ":" + cfg.ServerPort
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleListTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := s.generator.ListTemplates()
	if err != nil {
		s.sendError(w, "Failed to list templates", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TemplatesResponse{Templates: templates})
}

func (s *Server) handleCreateMeme(w http.ResponseWriter, r *http.Request) {
	var req MemeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var outputPath string
	var err error

	if req.ImageURL != "" {
		outputPath, err = s.generator.CreateMemeFromURL(req.ImageURL, req.TopText, req.BottomText)
	} else if req.Template != "" {
		outputPath, err = s.generator.CreateMeme(req.Template, req.TopText, req.BottomText)
	} else {
		s.sendError(w, "Either template or image_url must be provided", http.StatusBadRequest)
		return
	}

	if err != nil {
		s.sendError(w, fmt.Sprintf("Failed to generate meme: %v", err), http.StatusInternalServerError)
		return
	}

	// Read the generated image
	imageData, err := os.ReadFile(outputPath)
	if err != nil {
		s.sendError(w, "Failed to read generated image", http.StatusInternalServerError)
		return
	}

	// Clean up the file
	defer os.Remove(outputPath)

	// Send the image
	w.Header().Set("Content-Type", "image/png")
	w.Write(imageData)
}

func (s *Server) sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
