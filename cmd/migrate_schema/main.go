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
		dbURL = "postgresql://postgres.gfzrowkthoqnubfxkhnt:M7P25z3P0wPE0VvD@aws-1-ap-southeast-1.pooler.supabase.com:5432/postgres" // fallback
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	queries := []string{
		"ALTER TABLE community_posts ADD COLUMN IF NOT EXISTS encrypted_content TEXT;",
		"ALTER TABLE community_posts ADD COLUMN IF NOT EXISTS content_type TEXT DEFAULT 'text';",
		"ALTER TABLE community_posts ADD COLUMN IF NOT EXISTS image_url TEXT;",
		"ALTER TABLE community_posts ADD COLUMN IF NOT EXISTS is_edited BOOLEAN DEFAULT FALSE;",
		"ALTER TABLE community_posts ADD COLUMN IF NOT EXISTS is_removed BOOLEAN DEFAULT FALSE;",
		// Optional: Drop 'content' if it conflicts or is unused, but safer to keep for now
		// "ALTER TABLE community_posts DROP COLUMN IF EXISTS content;",
	}

	for _, q := range queries {
		fmt.Printf("Executing: %s\n", q)
		if _, err := db.Exec(q); err != nil {
			log.Printf("Error: %v\n", err)
		} else {
			fmt.Println("Success.")
		}
	}
	fmt.Println("Schema update complete.")
}
