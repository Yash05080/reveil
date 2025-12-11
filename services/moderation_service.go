package services

import (
	"context"
	"fmt"
	"strings"

	"reveil-api/models"
)

// ModerationService handles content moderation
type ModerationService struct {
	// keywords []string // Removed in favor of moderation_list.go
}

// NewModerationService creates a new ModerationService
func NewModerationService() *ModerationService {
	// Convert slice to map for O(1) lookups?
	// Actually, strictly matching "phrases" requires substrings check, not exact word match.
	// The previous implementation likely looped through options.
	return &ModerationService{
		// We'll use the package level variable BlockedPhrases in CheckPost
	}
}

// CheckPost performs synchronous checks on post content
// Returns: isFlagged (bool), reason (FlagReason), error
// CheckPost performs synchronous checks on post content
// Returns: isFlagged (bool), reason (FlagReason), error
func (s *ModerationService) CheckPost(ctx context.Context, content string) (*models.ModerationCheckResult, error) {
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
