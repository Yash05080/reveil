package services

import (
	"context"
	"strings"

	"reveil-api/models"
)

// ModerationService handles content moderation
type ModerationService struct {
	// In the future, we might have DB access or external API clients here
}

// NewModerationService creates a new ModerationService
func NewModerationService() *ModerationService {
	return &ModerationService{}
}

// CheckPost performs synchronous checks on post content
// Returns: isFlagged (bool), reason (FlagReason), error
func (s *ModerationService) CheckPost(ctx context.Context, content string) (models.ModerationCheckResult, error) {
	// 1. Keyword Check (Light Moderation)
	if match, reason := s.containsRestrictedKeywords(content); match {
		return models.ModerationCheckResult{
			IsFlagged:       true,
			FlagReason:      reason,
			SeverityLevel:   5, // Max severity for restricted keywords
			ConfidenceScore: 1.0,
			ShouldBlock:     false, // We flag but allow (or allow but hide? Plan said flag silently)
			// Actually, let's follow the plan: "Flag silently for now"
		}, nil
	}

	return models.ModerationCheckResult{
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
