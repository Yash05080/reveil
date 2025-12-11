package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("SUPABASE_URL")
	if dbURL == "" {
		dbURL = "postgresql://postgres.gfzrowkthoqnubfxkhnt:M7P25z3P0wPE0VvD@aws-1-ap-southeast-1.pooler.supabase.com:5432/postgres"
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	communityID := "aaaaaaa1-aaaa-aaaa-aaaa-aaaaaaaaaaaa" // From logs

	// Check for NULL encrypted_content
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM community_posts WHERE community_id = $1 AND encrypted_content IS NULL", communityID).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Rows with NULL encrypted_content: %d\n", count)

	// Check total rows
	var total int
	err = db.QueryRow("SELECT COUNT(*) FROM community_posts WHERE community_id = $1", communityID).Scan(&total)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total rows: %d\n", total)
}
