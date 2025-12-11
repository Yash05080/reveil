package models

import "github.com/google/uuid"

// PostLike mirrors public.post_likes
type PostLike struct {
	PostID uuid.UUID `db:"post_id"`
	UserID uuid.UUID `db:"user_id"`
}

// ToggleLikeResponse
type ToggleLikeResponse struct {
	Liked    bool `json:"liked"`     // True if liked, False if unliked
	NewCount int  `json:"new_count"` // Updated count to show immediately
}
