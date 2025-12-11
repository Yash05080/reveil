package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"reveil-api/config"
	"reveil-api/utils"
)

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	Sub         string `json:"sub"`
	CommunityID string `json:"community_id"`
	Email       string `json:"email"`
	Iat         int64  `json:"iat"`
	Exp         int64  `json:"exp"`
	Role        string `json:"role,omitempty"`
}

// Auth validates JWT tokens and extracts user information
func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.ErrorResponseWithCode(w, http.StatusUnauthorized,
					"Authorization header required", config.ErrorAuthentication)
				return
			}

			// Extract token from "Bearer <token>" format
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.ErrorResponseWithCode(w, http.StatusUnauthorized,
					"Invalid authorization header format", config.ErrorAuthentication)
				return
			}

			token := parts[1]

			// Validate and parse JWT token
			claims, err := validateJWT(token, jwtSecret)
			if err != nil {
				utils.ErrorResponseWithCode(w, http.StatusUnauthorized,
					fmt.Sprintf("Invalid token: %v", err), config.ErrorAuthentication)
				return
			}

			// Check token expiration
			if claims.Exp < time.Now().Unix() {
				utils.ErrorResponseWithCode(w, http.StatusUnauthorized,
					"Token expired", config.ErrorAuthentication)
				return
			}

			// Add user information to request context
			ctx := context.WithValue(r.Context(), "user_id", claims.Sub)
			ctx = context.WithValue(ctx, "community_id", claims.CommunityID)
			ctx = context.WithValue(ctx, "user_email", claims.Email)
			ctx = context.WithValue(ctx, "user_role", claims.Role)

			// Continue to next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// validateJWT validates a JWT token and returns the claims
func validateJWT(token, secret string) (*JWTClaims, error) {
	// Split token into parts
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	fmt.Printf("DEBUG: Validating Token Parts:\nHeader: %s\nPayload: %s\nSignature: %s\n", parts[0], parts[1], parts[2])

	// Decode payload (handle both RawURL and standard URL encoding)
	payload, err := parseBase64(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid payload encoding: %v", err)
	}

	signature, err := parseBase64(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid signature encoding: %v", err)
	}

	// Verify signature
	expectedSignature := generateSignature(parts[0]+"."+parts[1], secret)
	if !hmac.Equal(signature, expectedSignature) {
		return nil, fmt.Errorf("invalid signature")
	}

	// Parse claims
	var claims JWTClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("invalid claims format")
	}

	return &claims, nil
}

// generateSignature generates HMAC-SHA256 signature for JWT
func generateSignature(message, secret string) []byte {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return h.Sum(nil)
}

// RequireAuth is a middleware that requires authentication
func RequireAuth(jwtSecret string) func(http.Handler) http.Handler {
	return Auth(jwtSecret)
}

// OptionalAuth is a middleware that optionally parses auth but doesn't require it
func OptionalAuth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				// If auth header exists, validate it
				Auth(jwtSecret)(next).ServeHTTP(w, r)
			} else {
				// No auth header, continue without user context
				next.ServeHTTP(w, r)
			}
		})
	}
}
