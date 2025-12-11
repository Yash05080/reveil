package services

import (
	"context"
	"reveil-api/models"
	"testing"
)

func TestModerationService_CheckPost(t *testing.T) {
	service := NewModerationService()
	ctx := context.Background()

	tests := []struct {
		name         string
		content      string
		wantFlagged  bool
		wantReason   models.FlagReason
		wantSeverity int
	}{
		{
			name:        "Safe Content",
			content:     "Hello world, expecting a great day!",
			wantFlagged: false,
		},
		{
			name:         "Suicidal Ideation",
			content:      "I just want to kill myself",
			wantFlagged:  true,
			wantReason:   models.FlagReasonSuicidalIdeation,
			wantSeverity: 5,
		},
		{
			name:         "Self Harm",
			content:      "Thinking about self harm today",
			wantFlagged:  true,
			wantReason:   models.FlagReasonSelfHarm,
			wantSeverity: 5,
		},
		{
			name:         "Case Insensitive",
			content:      "I WANT TO DIE",
			wantFlagged:  true,
			wantReason:   models.FlagReasonSuicidalIdeation, // Matches "want to die"
			wantSeverity: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := service.CheckPost(ctx, tt.content)
			if got.IsFlagged != tt.wantFlagged {
				t.Errorf("CheckPost() isFlagged = %v, want %v", got.IsFlagged, tt.wantFlagged)
			}
			if tt.wantFlagged {
				if got.FlagReason != tt.wantReason {
					t.Errorf("CheckPost() flagReason = %v, want %v", got.FlagReason, tt.wantReason)
				}
				if got.SeverityLevel != tt.wantSeverity {
					t.Errorf("CheckPost() severityLevel = %v, want %v", got.SeverityLevel, tt.wantSeverity)
				}
			}
		})
	}
}
