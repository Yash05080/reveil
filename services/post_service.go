package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"reveil-api/models"

	"github.com/google/uuid"
)

type PostService struct {
	db                *sql.DB
	enc               *EncryptionService
	mod               *ModerationService
	sse               *SSEService
	AsyncAnalysisChan chan<- uuid.UUID
}

// NewPostService creates a new PostService
func NewPostService(db *sql.DB, enc *EncryptionService, mod *ModerationService, sse *SSEService, asyncChan chan<- uuid.UUID) *PostService {
	return &PostService{
		db:                db,
		enc:               enc,
		mod:               mod,
		sse:               sse,
		AsyncAnalysisChan: asyncChan,
	}
}

// CreatePost stores an encrypted post for a community
func (s *PostService) CreatePost(ctx context.Context, communityID, userID uuid.UUID, req models.CreatePostRequest) (*models.PostResponse, error) {
	log.Println("DEBUG: Entering CreatePost")
	// 1. Light Moderation Check (Content AND Title)
	// Checking title is important too
	modResult, _ := s.mod.CheckPost(ctx, req.Title+" "+req.Content)
	// We ignore error for now or log it, as per MVP.

	encryptedTitle, err := s.enc.EncryptContent(communityID, req.Title)
	if err != nil {
		return nil, fmt.Errorf("encrypt title failed: %w", err)
	}

	encryptedContent, err := s.enc.EncryptContent(communityID, req.Content)
	if err != nil {
		log.Printf("DEBUG: EncryptContent Error: %v\n", err)
		return nil, fmt.Errorf("encrypt content failed: %w", err)
	}

	var post models.CommunityPost
	err = s.db.QueryRowContext(ctx, `
		INSERT INTO public.community_posts (
			community_id, user_id, encrypted_title, encrypted_content, content_type, image_url,
			like_count, comment_count, created_at, updated_at, is_edited, is_removed
		) VALUES ($1, $2, $3, $4, $5, $6, 0, 0, NOW(), NOW(), FALSE, FALSE)
		RETURNING id, community_id, user_id, encrypted_title, encrypted_content, content_type, image_url,
		          like_count, comment_count, created_at, updated_at, is_edited, is_removed
	`, communityID, userID, encryptedTitle, encryptedContent, req.ContentType, req.ImageURL).Scan(
		&post.ID,
		&post.CommunityID,
		&post.UserID,
		&post.EncryptedTitle,
		&post.EncryptedContent,
		&post.ContentType,
		&post.ImageURL,
		&post.LikeCount,
		&post.CommentCount,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.IsEdited,
		&post.IsRemoved,
	)
	if err != nil {
		log.Printf("DEBUG: CreatePost DB Error: %v\n", err)
		return nil, fmt.Errorf("insert post failed: %w", err)
	}

	// 2. Insert Moderation Flag if needed
	if modResult.IsFlagged {
		_, err = s.db.ExecContext(ctx, `
			INSERT INTO public.moderation_flags (
				post_id, flag_reason, checked_by, action_taken,
				severity_level, confidence_score, notified_moderator
			) VALUES ($1, $2, $3, $4, $5, $6, FALSE)
		`, post.ID, modResult.FlagReason, models.CheckedByLightModel, models.ActionMarked,
			modResult.SeverityLevel, modResult.ConfidenceScore)
		if err != nil {
			fmt.Printf("Failed to insert moderation flag: %v\n", err)
		}
	}

	// 3. Broadcast Event via SSE
	resp, err := s.mapToResponse(ctx, post)
	if err == nil {
		// Only broadcast if we successfully mapped to response (decrypted)
		s.sse.BroadcastPostCreated(communityID, post.ID, resp)
	}

	// 4. Trigger Heavy Analysis (Async)
	if s.AsyncAnalysisChan != nil {
		select {
		case s.AsyncAnalysisChan <- post.ID:
			log.Println("DEBUG: Queued post for heavy analysis")
		default:
			log.Println("WARN: Analysis queue full, skipping heavy moderation")
		}
	}

	return resp, nil
}

// UpdatePost updates a post's content and image
func (s *PostService) UpdatePost(ctx context.Context, postID, userID uuid.UUID, req models.UpdatePostRequest) (*models.PostResponse, error) {
	// 1. Fetch existing post to get community_id (needed for encryption key) and verify ownership
	var post models.CommunityPost
	err := s.db.QueryRowContext(ctx, `
		SELECT id, community_id, user_id, encrypted_title, encrypted_content, content_type, image_url,
		       like_count, comment_count, created_at, updated_at, is_edited, is_removed
		FROM public.community_posts
		WHERE id = $1
	`, postID).Scan(
		&post.ID, &post.CommunityID, &post.UserID, &post.EncryptedTitle, &post.EncryptedContent, &post.ContentType, &post.ImageURL,
		&post.LikeCount, &post.CommentCount, &post.CreatedAt, &post.UpdatedAt, &post.IsEdited, &post.IsRemoved,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("post not found")
	}
	if err != nil {
		return nil, fmt.Errorf("fetch post failed: %w", err)
	}

	// 2. verify ownership and not removed
	if post.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}
	if post.IsRemoved {
		return nil, fmt.Errorf("post is removed")
	}

	// 2.5 New Moderation Check on Edit
	checkContent := req.Content
	if req.Title != "" {
		checkContent = req.Title + " " + req.Content
	}
	modResult, _ := s.mod.CheckPost(ctx, checkContent)

	// 3. Encrypt new content
	encryptedTitle := post.EncryptedTitle // Default to existing
	if req.Title != "" {
		encryptedTitle, err = s.enc.EncryptContent(post.CommunityID, req.Title)
		if err != nil {
			return nil, fmt.Errorf("encrypt title failed: %w", err)
		}
	}

	newEncryptedContent, err := s.enc.EncryptContent(post.CommunityID, req.Content)
	if err != nil {
		return nil, fmt.Errorf("encrypt content failed: %w", err)
	}

	// 4. Update in DB
	err = s.db.QueryRowContext(ctx, `
		UPDATE public.community_posts
		SET encrypted_title = $1, encrypted_content = $2, image_url = $3, is_edited = TRUE, updated_at = NOW()
		WHERE id = $4
		RETURNING id, community_id, user_id, encrypted_title, encrypted_content, content_type, image_url,
		          like_count, comment_count, created_at, updated_at, is_edited, is_removed
	`, encryptedTitle, newEncryptedContent, req.ImageURL, postID).Scan(
		&post.ID, &post.CommunityID, &post.UserID, &post.EncryptedTitle, &post.EncryptedContent, &post.ContentType, &post.ImageURL,
		&post.LikeCount, &post.CommentCount, &post.CreatedAt, &post.UpdatedAt, &post.IsEdited, &post.IsRemoved,
	)
	if err != nil {
		return nil, fmt.Errorf("update post failed: %w", err)
	}

	// 5. Apply Moderation Flags (Edit could introduce toxicity)
	if modResult.IsFlagged {
		_, _ = s.db.ExecContext(ctx, `
			INSERT INTO public.moderation_flags (
				post_id, flag_reason, checked_by, action_taken,
				severity_level, confidence_score, notified_moderator
			) VALUES ($1, $2, $3, $4, $5, $6, FALSE)
		`, post.ID, modResult.FlagReason, models.CheckedByLightModel, models.ActionMarked,
			modResult.SeverityLevel, modResult.ConfidenceScore)
	}

	// 6. Trigger Heavy Analysis for Edit
	if s.AsyncAnalysisChan != nil {
		select {
		case s.AsyncAnalysisChan <- post.ID:
			// Queued
		default:
			// Queue full
		}
	}

	return s.mapToResponse(ctx, post)
}

// ReportPost allows a user to flag a post
func (s *PostService) ReportPost(ctx context.Context, postID, userID uuid.UUID, reason string) error {
	// 1. Verify Post Exists
	var exists bool
	err := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM public.community_posts WHERE id=$1)", postID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("post not found")
	}

	// 2. Insert User Report
	// We prefix reason with "Report: " to distinguish
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO public.moderation_flags (
			post_id, flag_reason, checked_by, action_taken,
			severity_level, confidence_score, notified_moderator
		) VALUES ($1, $2, 'user_report', 'reported', 1, 1.0, FALSE)
	`, postID, "Report: "+reason)

	if err != nil {
		return fmt.Errorf("failed to report post: %w", err)
	}
	return nil
}

// DeletePost performs a soft delete
func (s *PostService) DeletePost(ctx context.Context, postID, userID uuid.UUID) error {
	result, err := s.db.ExecContext(ctx, `
		UPDATE public.community_posts
		SET is_removed = TRUE, updated_at = NOW()
		WHERE id = $1 AND user_id = $2 AND is_removed = FALSE
	`, postID, userID)
	if err != nil {
		return fmt.Errorf("delete post failed: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("post not found or unauthorized")
	}
	return nil
}

// ListPosts returns posts for a community with simple pagination
func (s *PostService) ListPosts(ctx context.Context, communityID uuid.UUID, q models.ListPostsQuery) ([]*models.PostResponse, error) {
	if q.Limit <= 0 || q.Limit > 50 {
		q.Limit = 20
	}

	args := []interface{}{communityID}
	query := `
		SELECT id, community_id, user_id, encrypted_title, encrypted_content, content_type, image_url,
		       like_count, comment_count, created_at, updated_at, is_edited, is_removed
		FROM public.community_posts
		WHERE community_id = $1
	`
	if q.Before != nil {
		query += fmt.Sprintf(" AND created_at < $%d", len(args)+1)
		args = append(args, *q.Before)
	}
	if q.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", len(args)+1)
		args = append(args, *q.UserID)
	}
	if q.ContentType != nil {
		query += fmt.Sprintf(" AND content_type = $%d", len(args)+1)
		args = append(args, *q.ContentType)
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d", len(args)+1)
	args = append(args, q.Limit)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("DEBUG: ListPosts DB Error: %v\n", err)
		return nil, fmt.Errorf("list posts failed: %w", err)
	}
	defer rows.Close()

	var result []*models.PostResponse
	for rows.Next() {
		var p models.CommunityPost
		if err := rows.Scan(
			&p.ID,
			&p.CommunityID,
			&p.UserID,
			&p.EncryptedTitle,
			&p.EncryptedContent,
			&p.ContentType,
			&p.ImageURL,
			&p.LikeCount,
			&p.CommentCount,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.IsEdited,
			&p.IsRemoved,
		); err != nil {
			return nil, err
		}
		resp, err := s.mapToResponse(ctx, p)
		if err != nil {
			return nil, err
		}
		result = append(result, resp)
	}
	return result, nil
}

func (s *PostService) mapToResponse(ctx context.Context, p models.CommunityPost) (*models.PostResponse, error) {
	title, err := s.enc.DecryptContent(p.CommunityID, p.EncryptedTitle)
	if err != nil {
		// Fallback for empty/legacy titles?
		title = ""
	}

	content, err := s.enc.DecryptContent(p.CommunityID, p.EncryptedContent)
	if err != nil {
		return nil, fmt.Errorf("decrypt failed: %w", err)
	}

	// Fetch highest severity flag if exists
	var modStatus *models.ModerationStatus
	var flagReason string
	var severityLevel int
	err = s.db.QueryRowContext(ctx, `
		SELECT flag_reason, severity_level 
		FROM public.moderation_flags 
		WHERE post_id = $1 
		ORDER BY severity_level DESC 
		LIMIT 1
	`, p.ID).Scan(&flagReason, &severityLevel)

	if err == nil {
		modStatus = &models.ModerationStatus{
			IsFlagged:     true,
			FlagReason:    flagReason,
			SeverityLevel: severityLevel,
		}
	} else if err != sql.ErrNoRows {
		// Log error but don't fail the request
		log.Printf("WARN: Failed to fetch flags for post %s: %v", p.ID, err)
	}

	return &models.PostResponse{
		ID:           p.ID,
		CommunityID:  p.CommunityID,
		UserID:       p.UserID,
		Title:        title,
		Content:      content,
		ContentType:  p.ContentType,
		ImageURL:     p.ImageURL,
		LikeCount:    p.LikeCount,
		CommentCount: p.CommentCount,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
		IsEdited:     p.IsEdited,
		IsRemoved:    p.IsRemoved,
		Moderation:   modStatus,
	}, nil
}
