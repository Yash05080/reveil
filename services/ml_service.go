package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MLService struct {
	BaseURL string
	client  *http.Client
}

type MLAnalysisResponse struct {
	IsAbusive       bool      `json:"is_abusive"`
	ConfidenceScore float64   `json:"confidence_score"`
	Flags           []FlagTag `json:"flags"`
}

type FlagTag struct {
	Label string  `json:"label"`
	Score float64 `json:"score"`
}

func NewMLService(baseURL string) *MLService {
	// Default to localhost:5001 if empty
	if baseURL == "" {
		baseURL = "http://localhost:5001"
	}
	return &MLService{
		BaseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// AnalyzeContent sends content to Python sidecar
func (s *MLService) AnalyzeContent(content string) (*MLAnalysisResponse, error) {
	payload := map[string]string{"content": content}
	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", s.BaseURL+"/analyze", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ml service unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ml service error: %d", resp.StatusCode)
	}

	var result MLAnalysisResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("invalid response json: %w", err)
	}

	return &result, nil
}
