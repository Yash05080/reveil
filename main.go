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
	var postService *services.PostService

	// Heavy Moderation Components
	mlService := services.NewMLService("http://localhost:5001") // Python sidecar
	// We create the worker first because it *owns* the channel or we create PostService?
	// Worker constructor creates the channel. Ideally PostService just takes `chan<-` interface.
	// But Cyclic dependency if Worker needs PostService to fetch content, and PostService needs Worker channel.
	// Circular dependency: HeavyWorker -> PostService -> AsyncChan.
	// Resolution:
	// 1. Create PostService with NIL channel first? No, constructor requires it.
	// 2. Create Channel explicitly here.

	// Create Worker (needs PostService to fetch content)
	// IMPORTANT: Worker needs full PostService? Or just DB and Encryption?
	// In heavy_worker.go, I defined `postService *PostService`.
	// So I can pass the *already created* postService to the Worker.
	// BUT the worker creates its own channel in `NewHeavyWorker` right now.
	// I should modify `NewHeavyWorker` to ACCEPT the channel or I should update `main.go` logic.
	// Let's modify `main.go` to construct properly.

	// I'll manually set the queue in main for simplicity or refactor worker.
	// Refactoring worker is better design, but editing main is faster if I just do:
	// worker := services.NewHeavyWorker(postService, mlService, dbConn)
	// worker.TaskQueue = asyncChan -- But TaskQueue is field.

	// Better Loop:
	// 1. Create components.
	// 2. Create PostService.
	// 3. Create Worker (passing PostService).
	// 4. Worker *provides* the channel?
	// Let's look at `NewHeavyWorker` again. It returns `*HeavyWorker`. It initializes `TaskQueue`.
	// So:
	// worker := services.NewHeavyWorker(nil, mlService, dbConn) -- PostService is nil?
	// postService := services.NewPostService(..., worker.TaskQueue)
	// worker.postService = postService -- set it back.

	// Let's go with Circular Dependency resolution by setting field after.

	// Step 1: Init ML

	// Step 2: Init Worker (with NIL postService for now)
	worker := services.NewHeavyWorker(nil, mlService, dbConn)

	// Step 3: Init PostService (using worker.TaskQueue)
	// Step 3: Init PostService (using worker.TaskQueue)
	// Redeclare or assign? Up above at line 68 I declared it?
	// Wait, line 68: postService := ...
	// And I am redeclaring it here.
	// I should remove the earlier declaration (line 68) or assign to it.
	// But line 68 used `asyncChan` which I am not using anymore (using worker.TaskQueue).
	// So I should remove the block from lines 65-68.

	// Actually, I can just use `=` here if I keep line 68, but line 68 is WRONG because it uses `asyncChan` which is created locally. I want to use `worker.TaskQueue`.
	// I should DELETE lines 65-68 as they are superseded by this block.

	// However, since I am editing line 102 only, I will change `:=` to `=` IF `postService` was already declared.
	// But `postService` is declared at line 68.
	// I should remove the previous declaration in a separate edit or just overwrite it here.

	// Better: Remove lines 65-68 completely in one edit, and verify this line is correct.

	// Let's do a multi-edit? No, sequential.
	// I will replace this line to use `=` assuming I keep the top one? NO, the top one uses `asyncChan` which is disconnected from worker.
	// The worker has its own queue.

	// I MUST remove lines 65-68.

	postService = services.NewPostService(dbConn, encService, modService, sseService, worker.TaskQueue)

	// Step 4: Inject PostService back into Worker
	worker.SetPostService(postService) // Need to add this setter or just assign if public field.
	// `postService` field in HeavyWorker is private (lower case). I need to make it public or add Setter.
	// I will add a Setter to `HeavyWorker` efficiently in next step. For now I must update main.go logic.

	// WAIT. Modifying `services/heavy_worker.go` to make `PostService` public or add setter is safer.
	// Let's assume I will add `SetPostService`.

	// Post handler and routes (auth required)
	postHandler := handlers.NewPostHandler(postService, sseService)
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.Auth(cfg.JWTSecret))
	postHandler.RegisterPostRoutes(api, validator)

	// Wrap with middleware
	handler := middleware.CORS(
		middleware.Logger(router, logger),
	)

	// Start Worker
	worker.Start()

	// Start server
	logger.Info("Server starting", "port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
