package services

import (
	"context"
	"database/sql"
	"log"

	"reveil-api/models"

	"github.com/google/uuid"
)

type LikeService struct {
	db *sql.DB
}

func NewLikeService(db *sql.DB) *LikeService {
	return &LikeService{db: db}
}

// TogglePostLike adds or removes a like
func (s *LikeService) TogglePostLike(ctx context.Context, postID, userID uuid.UUID) (*models.ToggleLikeResponse, error) {
	// Check if already liked
	var exists bool
	err := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM public.post_likes WHERE post_id=$1 AND user_id=$2)", postID, userID).Scan(&exists)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var liked bool
	if exists {
		// Unlike
		_, err = tx.ExecContext(ctx, "DELETE FROM public.post_likes WHERE post_id=$1 AND user_id=$2", postID, userID)
		if err != nil {
			return nil, err
		}
		_, err = tx.ExecContext(ctx, "UPDATE public.community_posts SET like_count = like_count - 1 WHERE id=$1", postID)
		liked = false
	} else {
		// Like
		_, err = tx.ExecContext(ctx, "INSERT INTO public.post_likes (post_id, user_id) VALUES ($1, $2)", postID, userID)
		if err != nil {
			return nil, err
		}
		_, err = tx.ExecContext(ctx, "UPDATE public.community_posts SET like_count = like_count + 1 WHERE id=$1", postID)
		liked = true
	}

	if err != nil {
		return nil, err
	}

	// Fetch new count
	var newCount int
	err = tx.QueryRowContext(ctx, "SELECT like_count FROM public.community_posts WHERE id=$1", postID).Scan(&newCount)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	log.Printf("User %s toggled like on Post %s. New State: %v", userID, postID, liked)

	return &models.ToggleLikeResponse{
		Liked:    liked,
		NewCount: newCount,
	}, nil
}
