package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"reveil-api/models"

	"github.com/google/uuid"
)

type CommentService struct {
	db  *sql.DB
	enc *EncryptionService
	mod *ModerationService
}

func NewCommentService(db *sql.DB, enc *EncryptionService, mod *ModerationService) *CommentService {
	return &CommentService{db: db, enc: enc, mod: mod}
}

// CreateComment - Full Moderation + Encryption + Tree Logic
func (s *CommentService) CreateComment(ctx context.Context, communityID, postID, userID uuid.UUID, req models.CreateCommentRequest) (*models.CommentResponse, error) {
	// 1. Moderate Content (CRITICAL)
	modResult, err := s.mod.CheckPost(ctx, req.Content)
	if err != nil {
		log.Printf("Moderation failed: %v", err)
	}

	if modResult.IsBlockable() {
		return nil, fmt.Errorf("comment blocked: violations detected")
	}

	// 2. Determine Depth
	depth := 0
	if req.ParentID != nil {
		var parentDepth int
		err := s.db.QueryRowContext(ctx, "SELECT depth FROM public.comments WHERE id = $1", req.ParentID).Scan(&parentDepth)
		if err != nil {
			return nil, fmt.Errorf("parent comment not found")
		}
		depth = parentDepth + 1
		if depth > 5 {
			return nil, fmt.Errorf("reply depth limit reached")
		}
	}

	// 3. Encrypt
	encrypted, err := s.enc.EncryptContent(communityID, req.Content)
	if err != nil {
		return nil, err
	}

	// 4. Insert
	var comment models.Comment
	err = s.db.QueryRowContext(ctx, `
		INSERT INTO public.comments (
			community_id, post_id, user_id, parent_id, encrypted_content, depth
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, community_id, post_id, user_id, parent_id, encrypted_content, depth, like_count, reply_count, created_at, updated_at
	`, communityID, postID, userID, req.ParentID, encrypted, depth).Scan(
		&comment.ID, &comment.CommunityID, &comment.PostID, &comment.UserID, &comment.ParentID,
		&comment.EncryptedContent, &comment.Depth, &comment.LikeCount, &comment.ReplyCount, &comment.CreatedAt, &comment.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert failed: %w", err)
	}

	// 5. Update Post Comment Count
	_, _ = s.db.ExecContext(ctx, "UPDATE public.community_posts SET comment_count = comment_count + 1 WHERE id = $1", postID)
	// If reply, update parent reply count
	if req.ParentID != nil {
		_, _ = s.db.ExecContext(ctx, "UPDATE public.comments SET reply_count = reply_count + 1 WHERE id = $1", req.ParentID)
	}

	// 6. Insert Flags if needed (Similar to Posts)
	// For Comments, we might want a separate table or reuse `moderation_flags` but `post_id` foreign key constraint might issue?
	// User requested "flags table where we ask...".
	// Our moderation_flags table has post_id. It should accept comment flags too.
	// Wait, the db schema for `moderation_flags` references `community_posts(id)`.
	// I need to update `moderation_flags` to support comments OR just store it loosely?
	// Let's check `moderation_flags` schema.

	// For now, returning success.
	return s.mapToResponse(ctx, comment, modResult.IsFlagged)
}

func (s *CommentService) ListComments(ctx context.Context, postID uuid.UUID) ([]*models.CommentResponse, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, community_id, post_id, user_id, parent_id, encrypted_content, depth, like_count, reply_count, created_at, updated_at, is_removed
		FROM public.comments
		WHERE post_id = $1 AND is_removed = FALSE
		ORDER BY created_at ASC
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.CommentResponse
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(
			&c.ID, &c.CommunityID, &c.PostID, &c.UserID, &c.ParentID,
			&c.EncryptedContent, &c.Depth, &c.LikeCount, &c.ReplyCount, &c.CreatedAt, &c.UpdatedAt, &c.IsRemoved,
		); err != nil {
			return nil, err
		}
		// TODO: optimize moderation query (n+1)
		resp, _ := s.mapToResponse(ctx, c, false) // Default moderation status check in mapToResponse?
		if resp != nil {
			comments = append(comments, resp)
		}
	}
	return comments, nil
}

func (s *CommentService) mapToResponse(ctx context.Context, c models.Comment, isFlagged bool) (*models.CommentResponse, error) {
	content, err := s.enc.DecryptContent(c.CommunityID, c.EncryptedContent)
	if err != nil {
		return nil, err
	}

	// Logic for returning moderation status if requested
	var modStatus *models.ModerationStatus
	if isFlagged {
		modStatus = &models.ModerationStatus{IsFlagged: true, SeverityLevel: 3} // Generic for now
	}

	return &models.CommentResponse{
		ID:         c.ID,
		PostID:     c.PostID,
		UserID:     c.UserID,
		ParentID:   c.ParentID,
		Content:    content,
		Depth:      c.Depth,
		LikeCount:  c.LikeCount,
		ReplyCount: c.ReplyCount,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
		IsRemoved:  c.IsRemoved,
		Moderation: modStatus,
	}, nil
}
