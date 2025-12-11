package utils

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
)

// GenerateRandomKey generates a random key of specified length
func GenerateRandomKey(length int) (string, error) {
    bytes := make([]byte, length)
    if _, err := rand.Read(bytes); err != nil {
        return "", fmt.Errorf("failed to generate random key: %w", err)
    }
    return base64.StdEncoding.EncodeToString(bytes), nil
}

// GenerateNonce generates a random nonce for encryption
func GenerateNonce(length int) ([]byte, error) {
    nonce := make([]byte, length)
    if _, err := rand.Read(nonce); err != nil {
        return nil, fmt.Errorf("failed to generate nonce: %w", err)
    }
    return nonce, nil
}

// Base64Encode encodes bytes to base64 string
func Base64Encode(data []byte) string {
    return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode decodes base64 string to bytes
func Base64Decode(data string) ([]byte, error) {
    return base64.StdEncoding.DecodeString(data)
}
