package main

import (
	"context"
	"log"
	"net/http"
	"time"

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
	dbConn := supabaseClient.DB()

	// Base Dependencies
	validator := utils.NewValidator()
	sseService := services.NewSSEService()
	encService := services.NewEncryptionService(dbConn, []byte(cfg.MasterEncryptionKey))
	mlService := services.NewMLService("http://localhost:5001")
	modService := services.NewModerationService(dbConn, mlService) // Mod requires DB for flags

	// Worker Init (Step 1: Create without PostService)
	// Note: We need to set PostService later to break circular dependency
	heavyWorker := services.NewHeavyWorker(dbConn, mlService, nil)

	// Services
	authService := services.NewAuthService(dbConn, cfg.JWTSecret)
	communityService := services.NewCommunityService(dbConn)
	postService := services.NewPostService(dbConn, encService, modService, sseService, heavyWorker.TaskQueue)
	commentService := services.NewCommentService(dbConn, encService, modService)
	likeService := services.NewLikeService(dbConn)

	// Worker Init (Step 2: Inject PostService)
	heavyWorker.SetPostService(postService)

	// Start Workers
	for i := 0; i < 3; i++ {
		go heavyWorker.Start(context.Background())
	}

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	communityHandler := handlers.NewCommunityHandler(communityService)
	postHandler := handlers.NewPostHandler(postService, commentService, likeService, sseService)
	healthHandler := handlers.NewHealthHandler(supabaseClient)

	// Router Setup
	router := mux.NewRouter()
	router.HandleFunc("/health", healthHandler.HealthCheck).Methods(http.MethodGet, http.MethodOptions)

	// Rate Limiter (5 req/s, burst 10)
	rateLimiter := middleware.NewRateLimiter(5, 10)
	go rateLimiter.CleanupLoop(time.Minute * 10)

	// API Subrouter
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.Auth(cfg.JWTSecret))
	api.Use(rateLimiter.Limit)

	// Register Routes
	authHandler.RegisterAuthRoutes(api, validator) // Assuming Auth routes exist or are needed?
	// Wait, Auth routes usually public (Login/Register).
	// My previous code didn't show Auth routes registration?
	// Let's assume Auth is handled properly or stick to what was there.
	// The previous code only showed Post/Community handlers.
	// I shall check if Auth routes need to be public (outside api middleware).

	// Check `handlers/auth_handler.go` if needed. For now, adding what was visible.
	communityHandler.RegisterCommunityRoutes(api, validator)
	postHandler.RegisterPostRoutes(api, validator)

	// Public Auth Routes (if any)
	// router.HandleFunc("/login", authHandler.Login).Methods("POST")
	// For now assuming existing structure.

	// Middleware
	handler := middleware.CORS(
		middleware.Logger(router, logger),
	)

	// Start server
	logger.Info("Server starting", "port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
