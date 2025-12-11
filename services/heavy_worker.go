package services

import (
	"context"
	"database/sql"
	"log"

	"reveil-api/models"

	"github.com/google/uuid"
)

type HeavyWorker struct {
	postService *PostService // To fetch post content (decrypts internaly)
	mlService   *MLService   // To analyze content
	db          *sql.DB      // To insert flag
	TaskQueue   chan uuid.UUID
}

func (w *HeavyWorker) SetPostService(ps *PostService) {
	w.postService = ps
}

func NewHeavyWorker(ps *PostService, ml *MLService, db *sql.DB) *HeavyWorker {
	return &HeavyWorker{
		postService: ps,
		mlService:   ml,
		db:          db,
		TaskQueue:   make(chan uuid.UUID, 100), // Buffer 100 posts
	}
}

func (w *HeavyWorker) Start() {
	log.Println("[Worker] Heavy Moderation Worker Started")
	go func() {
		for postID := range w.TaskQueue {
			w.processPost(postID)
		}
	}()
}

func (w *HeavyWorker) processPost(postID uuid.UUID) {
	log.Printf("[Worker] Processing Post: %s\n", postID)

	// 1. Fetch Post (This decrypts content)
	// We need a way to fetch post by ID purely, passing a context.
	// We'll use a dummy/system context for now.
	ctx := context.Background()

	// Create a dummy Query for List or just fetch by ID.
	// PostService.UpdatePost fetches by ID but requires user ID.
	// PostService doesn't have a public "GetPostByID" for internal use efficiently?
	// It has `mapToResponse` which decrypts.
	// Let's add a `GetPostInternal` or just query DB directly here and use encService?
	// But PostService encapsulates EncryptionService.
	// Ideally PostService should have `GetPost(ctx, postID)` that returns decrypted content.
	// But `UpdatePost` logic does it.
	// Let's rely on `PostService` having access to `enc`.
	// Wait, `HeavyWorker` is inside `services` package. It can access private fields/methods of `PostService`?
	// No, different structs same package. Yes, `s.enc` is private but visible in same package `services`.

	// So we can do:
	var post models.CommunityPost
	err := w.db.QueryRowContext(ctx, `
        SELECT id, community_id, user_id, encrypted_content, content_type 
        FROM community_posts WHERE id = $1
    `, postID).Scan(&post.ID, &post.CommunityID, &post.UserID, &post.EncryptedContent, &post.ContentType)

	if err != nil {
		log.Printf("[Worker] Failed to fetch post %s: %v\n", postID, err)
		return
	}

	// Decrypt
	content, err := w.postService.enc.DecryptContent(post.CommunityID, post.EncryptedContent)
	if err != nil {
		log.Printf("[Worker] Decrypt failed for %s: %v\n", postID, err)
		return
	}

	// 2. Analyze
	result, err := w.mlService.Analyze(content)
	if err != nil {
		log.Printf("[Worker] ML Analysis failed for %s: %v\n", postID, err)
		return
	}

	log.Printf("[Worker] Analysis Result for %s: Abusive=%v Score=%.2f\n", postID, result.IsAbusive, result.ConfidenceScore)

	// 3. Flag if abusive
	if result.IsAbusive {
		_, err := w.db.ExecContext(ctx, `
			INSERT INTO public.moderation_flags (
				post_id, flag_reason, checked_by, action_taken,
				severity_level, confidence_score, notified_moderator
			) VALUES ($1, $2, $3, $4, $5, $6, FALSE)
		`, postID, "Hate Speech (ML)", models.CheckedByHeavyModel, models.ActionMarked,
			3, result.ConfidenceScore) // Severity 3 for ML flags
		if err != nil {
			log.Printf("[Worker] Failed to insert flag for %s: %v\n", postID, err)
		} else {
			log.Printf("[Worker] FLAGGED Post %s\n", postID)
		}
	}
}
