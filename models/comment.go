package models

import (
	"time"

	"github.com/google/uuid"
)

// Comment mirrors public.comments
type Comment struct {
	ID               uuid.UUID  `db:"id"`
	CommunityID      uuid.UUID  `db:"community_id"`
	PostID           uuid.UUID  `db:"post_id"`
	UserID           uuid.UUID  `db:"user_id"`
	ParentID         *uuid.UUID `db:"parent_id"` // Nullable
	EncryptedContent string     `db:"encrypted_content"`
	Depth            int        `db:"depth"`
	LikeCount        int        `db:"like_count"`
	ReplyCount       int        `db:"reply_count"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
	IsRemoved        bool       `db:"is_removed"`
}

// CreateCommentRequest payload
type CreateCommentRequest struct {
	Content  string     `json:"content" validate:"required,min=1,max=2000"`
	ParentID *uuid.UUID `json:"parent_id,omitempty"` // Optional for top-level
}

// CommentResponse for frontend
type CommentResponse struct {
	ID         uuid.UUID         `json:"id"`
	PostID     uuid.UUID         `json:"post_id"`
	UserID     uuid.UUID         `json:"user_id"`
	ParentID   *uuid.UUID        `json:"parent_id,omitempty"`
	Content    string            `json:"content"` // Decrypted
	Depth      int               `json:"depth"`
	LikeCount  int               `json:"like_count"`
	ReplyCount int               `json:"reply_count"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	IsRemoved  bool              `json:"is_removed"`
	Moderation *ModerationStatus `json:"moderation,omitempty"`
}
