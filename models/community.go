package models

import (
	"time"

	"github.com/google/uuid"
)

// Community represents a community in the database
type Community struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description,omitempty" db:"description"`
	CreatedBy   uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CommunityMember represents a user's membership in a community
type CommunityMember struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CommunityID uuid.UUID `json:"community_id" db:"community_id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Role        *string   `json:"role,omitempty" db:"role"`
	JoinedAt    time.Time `json:"joined_at" db:"joined_at"`
}

// EncryptionKey represents community-specific encryption keys
type EncryptionKey struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	CommunityID  uuid.UUID  `json:"community_id" db:"community_id"`
	EncryptedKey string     `json:"-" db:"encrypted_key"` // Never expose
	KeyVersion   int        `json:"key_version" db:"key_version"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	RotatedAt    *time.Time `json:"rotated_at,omitempty" db:"rotated_at"`
}

type CreateCommunityRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"max=500"`
}
