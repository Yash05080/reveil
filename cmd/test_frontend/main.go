package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	BaseURL = "http://localhost:8080/api"
	// Developer Token (Valid for testing)
	Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb21tdW5pdHlfaWQiOiJhYWFhYWFhMS1hYWFhLWFhYWEtYWFhYS1hYWFhYWFhYWFhYWEiLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJleHAiOjE3NjU1NDYyODEsImlhdCI6MTc2NTQ1OTg4MSwicm9sZSI6InVzZXIiLCJzdWIiOiIxMTExMTExMS0xMTExLTExMTEtMTExMS0xMTExMTExMTExMTEifQ.Q1gOjWFLd8IB3wskjA9eMywo9NYRt6hGp0fuDju2n4Y"
)

func main() {
	fmt.Println("=== REVEIL API FRONTEND TEST SCRIPT ===")

	// 1. Create Post with Title
	fmt.Println("\n[1] Creating Post with Title...")
	payload := map[string]interface{}{
		"title":        "Frontend Integration Test",
		"content":      "This post verifies the title field integration.",
		"content_type": "text",
		"image_url":    "https://example.com/image.png",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", BaseURL+"/communities/aaaaaaa1-aaaa-aaaa-aaaa-aaaaaaaaaaaa/posts", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %d\nResponse: %s\n", resp.StatusCode, string(respBody))

	if resp.StatusCode == 201 {
		fmt.Println("✅ Post Created Successfully with Title")
	} else {
		fmt.Println("❌ Failed to create Post")
	}

	fmt.Println("\n=== TEST COMPLETE ===")
}
