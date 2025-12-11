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
		// Hardcode from .env for convenience if env var not set in shell
		dbURL = "postgresql://postgres.gfzrowkthoqnubfxkhnt:M7P25z3P0wPE0VvD@aws-1-ap-southeast-1.pooler.supabase.com:5432/postgres"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(`
		SELECT column_name, data_type 
		FROM information_schema.columns 
		WHERE table_name = 'community_posts';
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Columns in community_posts:")
	found := false
	for rows.Next() {
		found = true
		var name, dtype string
		if err := rows.Scan(&name, &dtype); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("- %s (%s)\n", name, dtype)
	}
	if !found {
		fmt.Println("Table community_posts not found or has no columns.")
	}
}
