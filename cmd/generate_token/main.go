package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const jwtSecret = "b4fdbc1d-a81b-4181-b3db-42d6f94f40cc" // From .env

func main() {
	// Connect to DB to get a real Community ID
	dbURL := os.Getenv("SUPABASE_URL")
	if dbURL == "" {
		dbURL = "postgresql://postgres.gfzrowkthoqnubfxkhnt:M7P25z3P0wPE0VvD@aws-1-ap-southeast-1.pooler.supabase.com:5432/postgres"
	}
	db, err := sql.Open("postgres", dbURL)
	var communityID string
	if err == nil {
		err = db.QueryRow("SELECT id FROM communities LIMIT 1").Scan(&communityID)
		if err != nil {
			fmt.Println("Warning: Could not fetch real community ID, using random (expect FK errors if not exists)")
			communityID = uuid.New().String()
		}
		// db.Close() // Moved closing of DB to after fetching user ID
	} else {
		communityID = uuid.New().String()
	}

	var userID string
	if err == nil { // db is open
		err = db.QueryRow("SELECT id FROM users LIMIT 1").Scan(&userID)
		if err != nil {
			fmt.Println("Warning: Could not fetch real User ID, using random")
			userID = uuid.New().String()
		}
		db.Close() // Close DB after fetching both IDs
	} else {
		userID = uuid.New().String()
	}

	header := `{"alg":"HS256","typ":"JWT"}`
	claims := map[string]interface{}{
		"sub":          userID,
		"community_id": communityID,
		"email":        "test@example.com",
		"iat":          time.Now().Unix(),
		"exp":          time.Now().Add(24 * time.Hour).Unix(),
		"role":         "user",
	}

	claimsJSON, _ := json.Marshal(claims)

	encodedHeader := base64.RawURLEncoding.EncodeToString([]byte(header))
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsJSON)

	message := encodedHeader + "." + encodedClaims
	signature := computeHMAC(message, jwtSecret)

	token := message + "." + base64.RawURLEncoding.EncodeToString(signature)

	fmt.Printf("\n=== GENERATED DEV TOKEN ===\n")
	fmt.Printf("User ID: %s\n", userID)
	fmt.Printf("Community ID: %s\n", communityID)
	fmt.Printf("Token:\n%s\n", token)
	fmt.Printf("===========================\n\n")
	fmt.Printf("Example CURL:\n")
	fmt.Printf("curl -H \"Authorization: Bearer %s\" http://localhost:8080/api/communities/%s/posts\n", token, communityID)
}

func computeHMAC(message, secret string) []byte {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return h.Sum(nil)
}
