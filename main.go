package main

import (
	"log"
	"net/http"

	"reveil-api/config"
	"reveil-api/db"
	"reveil-api/handlers"
	"reveil-api/middleware"
	"reveil-api/services"
	"reveil-api/utils"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger := utils.NewLogger(cfg.LogLevel)
	logger.Info("Starting Reveil API server", "port", cfg.Port)

	// Initialize Supabase client
	supabaseClient, err := db.NewSupabaseClient(cfg.SupabaseURL, cfg.SupabaseServiceKey)
	if err != nil {
		log.Fatalf("Failed to connect to Supabase: %v", err)
	}
	defer supabaseClient.Close()

	// Test database connection
	if err := supabaseClient.Health(); err != nil {
		log.Fatalf("Database health check failed: %v", err)
	}
	logger.Info("Connected to Supabase successfully")

	// Base router
	router := mux.NewRouter()

	// Health handler (no auth)
	healthHandler := handlers.NewHealthHandler(supabaseClient)
	router.HandleFunc("/health", healthHandler.HealthCheck).Methods(http.MethodGet, http.MethodOptions)

	// Shared deps
	dbConn := supabaseClient.DB()
	validator := utils.NewValidator()
	encService := services.NewEncryptionService(dbConn, []byte(cfg.MasterEncryptionKey))
	modService := services.NewModerationService()
	sseService := services.NewSSEService()
	postService := services.NewPostService(dbConn, encService, modService, sseService)

	// Post handler and routes (auth required)
	postHandler := handlers.NewPostHandler(postService, sseService)
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.Auth(cfg.JWTSecret))
	postHandler.RegisterPostRoutes(api, validator)

	// Wrap with middleware
	handler := middleware.CORS(
		middleware.Logger(router, logger),
	)

	// Start server
	logger.Info("Server starting", "port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
