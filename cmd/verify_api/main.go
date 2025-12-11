package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	// 1. Generate Token
	fmt.Println(">>> Generating Token...")
	out, err := exec.Command("go", "run", "cmd/generate_token/main.go").Output()
	if err != nil {
		fmt.Printf("Failed to generate token: %v\n", err)
		os.Exit(1)
	}
	output := string(out)

	// Extract Token and Community ID
	token := extractValue(output, "Token:\n", "\n")
	communityID := extractValue(output, "Community ID: ", "\n")

	if token == "" || communityID == "" {
		fmt.Println("Failed to parse token output")
		fmt.Println(output)
		os.Exit(1)
	}
	fmt.Printf("Token: %s...\nCommunity: %s\n", token[:10], communityID)

	// 2. Test Post Creation (Auth Check)
	fmt.Println("\n>>> Testing Create Post (Auth Check)...")
	url := fmt.Sprintf("http://localhost:8080/api/communities/%s/posts", communityID)
	body := []byte(`{"content": "Automated test content", "content_type": "text"}`)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %s\nResponse: %s\n", resp.Status, string(respBody))

	if resp.StatusCode == 201 {
		fmt.Println("✅ Auth & Creation Success!")
	} else {
		fmt.Println("❌ Auth Failed")
		os.Exit(1)
	}

	// 3. Test Moderation (Flag Check)
	fmt.Println("\n>>> Testing Moderation (Self-harm keyword)...")
	bodyMod := []byte(`{"content": "I want to hurt myself", "content_type": "text"}`)
	reqMod, _ := http.NewRequest("POST", url, bytes.NewBuffer(bodyMod))
	reqMod.Header.Set("Authorization", "Bearer "+token)
	reqMod.Header.Set("Content-Type", "application/json")

	respMod, _ := client.Do(reqMod)
	if respMod.StatusCode == 201 {
		fmt.Println("✅ Moderation Request Accepted (will check logs for flag)")
	} else {
		fmt.Printf("❌ Moderation Request Failed: %s\n", respMod.Status)
	}

	// 4. Test List Posts
	fmt.Println("\n>>> Testing List Posts...")
	reqList, _ := http.NewRequest("GET", url, nil)
	reqList.Header.Set("Authorization", "Bearer "+token)
	respList, err := client.Do(reqList)
	if err != nil {
		fmt.Printf("List request failed: %v\n", err)
	} else {
		defer respList.Body.Close()
		bodyList, _ := io.ReadAll(respList.Body)
		if respList.StatusCode == 200 {
			fmt.Printf("✅ List Posts Success (200 OK)\nResponse: %s\n", string(bodyList))
		} else {
			fmt.Printf("❌ List Posts Failed: %s\n", respList.Status)
		}
	}
}

func extractValue(s, prefix, suffix string) string {
	start := strings.Index(s, prefix)
	if start == -1 {
		return ""
	}
	start += len(prefix)
	rest := s[start:]
	end := strings.Index(rest, suffix)
	if end == -1 {
		return rest
	}
	return rest[:end]
}
