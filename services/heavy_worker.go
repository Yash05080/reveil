package services

import (
	"context"
	"database/sql"
	"fmt"
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

func NewHeavyWorker(db *sql.DB, ml *MLService, ps *PostService) *HeavyWorker {
	return &HeavyWorker{
		db:          db,
		mlService:   ml,
		postService: ps,
		TaskQueue:   make(chan uuid.UUID, 100), // Buffer
	}
}

// SetPostService breaks circular dependency
func (w *HeavyWorker) SetPostService(ps *PostService) {
	w.postService = ps
}

func (w *HeavyWorker) Start(ctx context.Context) {
	log.Println("[Worker] Post Analysis Worker waiting for tasks...")
	for {
		select {
		case postID := <-w.TaskQueue:
			w.processPost(ctx, postID)
		case <-ctx.Done():
			return
		}
	}
}

func (w *HeavyWorker) processPost(ctx context.Context, postID uuid.UUID) {
	log.Printf("[Worker] Processing Post: %s\n", postID)

	// 1. Fetch Post (This decrypts content)
	// We need a way to fetch post by ID purely, passing a context.
	// We'll use a dummy/system context for now.

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
	decryptedContent, err := w.postService.enc.DecryptContent(post.CommunityID, post.EncryptedContent)
	if err != nil {
		log.Printf("[Worker] Decrypt failed for %s: %v\n", postID, err)
		return
	}

	// 2. Heavy ML Analysis (if not flagged by light model or always?)
	// Currently we always run ML as "second opinion" or "deeper check"
	result, err := w.mlService.AnalyzeContent(decryptedContent)
	if err != nil {
		log.Printf("[Worker] ML Analysis failed for %s: %v\n", postID, err)
		return
	}

	log.Printf("[Worker] Analysis Result for %s: Abusive=%v Score=%.2f\n", postID, result.IsAbusive, result.ConfidenceScore)

	// 3. Flag if abusive
	if result.IsAbusive {
		reason := "Toxic Content (ML)"
		if len(result.Flags) > 0 {
			// Pick the most severe label or join them?
			// Let's just take the first one or top one for the main reason.
			// Example: "threat" is worse than "toxic"
			// Simple approach: Use the label of the highest score (which `app.py` logic likely put in list, but we can iterate)
			// Actually `app.py` returns list.
			reason = fmt.Sprintf("ML: %s", result.Flags[0].Label)
		}

		_, err := w.db.ExecContext(ctx, `
			INSERT INTO public.moderation_flags (
				post_id, flag_reason, checked_by, action_taken,
				severity_level, confidence_score, notified_moderator
			) VALUES ($1, $2, $3, $4, $5, $6, FALSE)
		`, postID, reason, models.CheckedByHeavyModel, models.ActionMarked,
			3, result.ConfidenceScore) // Severity 3 for ML flags
		if err != nil {
			log.Printf("[Worker] Failed to insert flag for %s: %v\n", postID, err)
		} else {
			log.Printf("[Worker] FLAGGED Post %s (Reason: %s)\n", postID, reason)
		}
	}
}
