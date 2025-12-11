package models

import (
	"time"

	"github.com/google/uuid"
)

// CommunityPost mirrors public.community_posts
type CommunityPost struct {
	ID               uuid.UUID `db:"id"`
	CommunityID      uuid.UUID `db:"community_id"`
	UserID           uuid.UUID `db:"user_id"`
	EncryptedTitle   string    `db:"encrypted_title"` // Added
	EncryptedContent string    `db:"encrypted_content"`
	ContentType      string    `db:"content_type"`
	ImageURL         *string   `db:"image_url"`
	LikeCount        int       `db:"like_count"`
	CommentCount     int       `db:"comment_count"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
	IsEdited         bool      `db:"is_edited"`
	IsRemoved        bool      `db:"is_removed"`
}

// CreatePostRequest is the incoming payload from frontend
type CreatePostRequest struct {
	Title       string  `json:"title" validate:"required,min=3,max=100"` // Added
	Content     string  `json:"content" validate:"required,min=1,max=5000"`
	ContentType string  `json:"content_type" validate:"required,oneof=text image link"`
	ImageURL    *string `json:"image_url,omitempty" validate:"omitempty,url"`
}

// UpdatePostRequest is the payload for updating a post
type UpdatePostRequest struct {
	Title    string  `json:"title" validate:"omitempty,min=3,max=100"` // Added
	Content  string  `json:"content" validate:"required,min=1,max=5000"`
	ImageURL *string `json:"image_url,omitempty" validate:"omitempty,url"`
}

// PostResponse is what we return to the frontend
type PostResponse struct {
	ID           uuid.UUID         `json:"id"`
	CommunityID  uuid.UUID         `json:"community_id"`
	UserID       uuid.UUID         `json:"user_id"`
	Title        string            `json:"title"`   // Added
	Content      string            `json:"content"` // decrypted
	ContentType  string            `json:"content_type"`
	ImageURL     *string           `json:"image_url,omitempty"`
	LikeCount    int               `json:"like_count"`
	CommentCount int               `json:"comment_count"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	IsEdited     bool              `json:"is_edited"`
	IsRemoved    bool              `json:"is_removed"`
	Moderation   *ModerationStatus `json:"moderation,omitempty"`
}

// ModerationStatus contains public moderation info for the frontend
type ModerationStatus struct {
	IsFlagged     bool   `json:"is_flagged"`
	FlagReason    string `json:"flag_reason,omitempty"`
	SeverityLevel int    `json:"severity_level,omitempty"`
}

// ListPostsQuery holds query params for GET /posts
type ListPostsQuery struct {
	Limit       int
	Before      *time.Time
	UserID      *uuid.UUID
	ContentType *string
}
