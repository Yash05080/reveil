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

	res, err := db.Exec("DELETE FROM community_posts WHERE encrypted_content IS NULL")
	if err != nil {
		log.Fatal(err)
	}
	rows, _ := res.RowsAffected()
	fmt.Printf("Deleted %d bad rows.\n", rows)
}
