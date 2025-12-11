package models

import (
    "time"
    
    "github.com/google/uuid"
)

// FlagReason represents the reason for content flagging
type FlagReason string

const (
    FlagReasonSuicidalIdeation FlagReason = "suicidal_ideation"
    FlagReasonSelfHarm         FlagReason = "self_harm"
    FlagReasonProfanity        FlagReason = "profanity"
    FlagReasonEatingDisorder   FlagReason = "eating_disorder" 
    FlagReasonSubstanceAbuse   FlagReason = "substance_abuse"
    FlagReasonToxicity         FlagReason = "toxicity"
)

// ActionTaken represents the moderation action taken
type ActionTaken string

const (
    ActionRemoved ActionTaken = "removed"
    ActionHidden  ActionTaken = "hidden"
    ActionMarked  ActionTaken = "marked"
)

// CheckedBy represents which model performed the check
type CheckedBy string

const (
    CheckedByLightModel CheckedBy = "light_model"
    CheckedByHeavyModel CheckedBy = "heavy_model"
)

// ModerationFlag represents a content moderation flag in the database
type ModerationFlag struct {
    ID                uuid.UUID   `json:"id" db:"id"`
    PostID            uuid.UUID   `json:"post_id" db:"post_id"`
    CommentID         *uuid.UUID  `json:"comment_id,omitempty" db:"comment_id"` // For future comments
    FlagReason        FlagReason  `json:"flag_reason" db:"flag_reason"`
    CheckedBy         CheckedBy   `json:"checked_by" db:"checked_by"`
    ActionTaken       ActionTaken `json:"action_taken" db:"action_taken"`
    SeverityLevel     int         `json:"severity_level" db:"severity_level"`
    ConfidenceScore   float64     `json:"confidence_score" db:"confidence_score"`
    NotifiedModerator bool        `json:"notified_moderator" db:"notified_moderator"`
    FlaggedAt         time.Time   `json:"flagged_at" db:"flagged_at"`
}

// ModerationCheckResult represents the result of a moderation check
type ModerationCheckResult struct {
    IsFlagged       bool       `json:"is_flagged"`
    FlagReason      FlagReason `json:"flag_reason,omitempty"`
    SeverityLevel   int        `json:"severity_level"`     // 1-5 scale
    ConfidenceScore float64    `json:"confidence_score"`   // 0.0-1.0
    ShouldBlock     bool       `json:"should_block"`       // true if severity >= 4
}

// IsBlockable returns true if the result should block content creation
func (mcr ModerationCheckResult) IsBlockable() bool {
    return mcr.ShouldBlock || mcr.SeverityLevel >= 4
}
