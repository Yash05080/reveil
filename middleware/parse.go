package middleware

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// parseBase64 tries to decode string using RawURL first, then standard URL encoding
func parseBase64(s string) ([]byte, error) {
	if data, err := base64.RawURLEncoding.DecodeString(s); err == nil {
		return data, nil
	}
	if data, err := base64.URLEncoding.DecodeString(s); err == nil {
		return data, nil
	}
	// Try adding padding manually if needed
	if m := len(s) % 4; m != 0 {
		s += strings.Repeat("=", 4-m)
	}
	data, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		fmt.Printf("DEBUG: Failed to decode base64 string: %s (Error: %v)\n", s, err)
	}
	return data, err
}
