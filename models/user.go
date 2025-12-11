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
