package models

import (
    "time"
    
    "github.com/google/uuid"
)

// User represents a user in the system
type User struct {
    ID         uuid.UUID `json:"id" db:"id"`
    Username   string    `json:"username" db:"username"`
    Email      string    `json:"email" db:"email"`
    IsVerified bool      `json:"is_verified" db:"is_verified"`
    CreatedAt  time.Time `json:"created_at" db:"created_at"`
    UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
    // Password hash not exposed in JSON for security
    PasswordHash string `json:"-" db:"password_hash"`
}

// Comment represents a comment on a post
type Comment struct {
    ID               uuid.UUID  `json:"id" db:"id"`
    PostID           uuid.UUID  `json:"post_id" db:"post_id"`
    UserID           uuid.UUID  `json:"user_id" db:"user_id"`
    ParentCommentID  *uuid.UUID `json:"parent_comment_id,omitempty" db:"parent_comment_id"`
    EncryptedContent string     `json:"-" db:"encrypted_content"`
    DecryptedContent string     `json:"content,omitempty"`
    LikeCount        int        `json:"like_count" db:"like_count"`
    ReplyCount       int        `json:"reply_count" db:"reply_count"`
    IsEdited         bool       `json:"is_edited" db:"is_edited"`
    IsRemoved        bool       `json:"is_removed" db:"is_removed"`
    CreatedAt        time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}
