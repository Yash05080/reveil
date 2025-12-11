package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"reveil-api/models"
)

// ModerationService handles content moderation
type ModerationService struct {
	db *sql.DB
	ml *MLService
}

// NewModerationService creates a new ModerationService
func NewModerationService(db *sql.DB, ml *MLService) *ModerationService {
	return &ModerationService{
		db: db,
		ml: ml,
	}
}

// CheckPost performs synchronous checks on post content
// Returns: isFlagged (bool), reason (FlagReason), error
// CheckPost performs synchronous checks on post content
// Returns: isFlagged (bool), reason (FlagReason), error
func (s *ModerationService) CheckPost(ctx context.Context, content string) (*models.ModerationCheckResult, error) {
	// 1. Check for specific restrictred keywords (Suicide/Self-harm)
	if isRestricted, reason := s.containsRestrictedKeywords(content); isRestricted {
		return &models.ModerationCheckResult{
			IsFlagged:       true,
			FlagReason:      reason,
			SeverityLevel:   5,
			ConfidenceScore: 1.0,
			ShouldBlock:     true, // Typically we want to offer help resources, maybe block
		}, nil
	}

	lowerContent := strings.ToLower(content)
	for _, keyword := range BlockedPhrases {
		if strings.Contains(lowerContent, keyword) {
			return &models.ModerationCheckResult{
				IsFlagged:       true,
				FlagReason:      models.FlagReason(fmt.Sprintf("Contains blocked phrase: '%s'", keyword)),
				SeverityLevel:   5, // High severity
				ConfidenceScore: 1.0,
				ShouldBlock:     false,
			}, nil
		}
	}

	return &models.ModerationCheckResult{
		IsFlagged:       false,
		ConfidenceScore: 0.0,
	}, nil
}

// containsRestrictedKeywords checks for banned words
// Returns true and the reason if found
func (s *ModerationService) containsRestrictedKeywords(content string) (bool, models.FlagReason) {
	lowerContent := strings.ToLower(content)

	// Hardcoded keyword lists (MVP)
	// In production, these should come from DB or config
	suicideKeywords := []string{"kill myself", "suicide", "end my life", "want to die"}
	selfHarmKeywords := []string{"cut myself", "hurt myself", "self harm"}

	for _, word := range suicideKeywords {
		if strings.Contains(lowerContent, word) {
			return true, models.FlagReasonSuicidalIdeation
		}
	}

	for _, word := range selfHarmKeywords {
		if strings.Contains(lowerContent, word) {
			return true, models.FlagReasonSelfHarm
		}
	}

	return false, ""
}
