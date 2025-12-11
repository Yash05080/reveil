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

type AnalysisResponse struct {
	IsAbusive       bool    `json:"is_abusive"`
	ConfidenceScore float64 `json:"confidence_score"`
	RawLabel        string  `json:"raw_label"`
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

func (s *MLService) Analyze(content string) (*AnalysisResponse, error) {
	payload := map[string]string{"content": content}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := s.client.Post(s.BaseURL+"/analyze", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("ml service request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("ml service returned status: %s", resp.Status)
	}

	var result AnalysisResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode ml response: %w", err)
	}

	return &result, nil
}
