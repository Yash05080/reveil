package db

import (
    "database/sql"
    "fmt"
    "time"
    
    _ "github.com/lib/pq" // PostgreSQL driver
)

// SupabaseClient wraps the database connection
type SupabaseClient struct {
    db *sql.DB
}

// NewSupabaseClient creates a new Supabase client connection
func NewSupabaseClient(databaseURL, serviceKey string) (*SupabaseClient, error) {
    // Parse connection string
    // For Supabase, the connection string format is:
    // postgres://postgres:[password]@db.[project-id].supabase.co:5432/postgres
    connStr := fmt.Sprintf("%s?sslmode=require", databaseURL)
    
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(20)
    db.SetMaxIdleConns(10)
    db.SetConnMaxLifetime(time.Hour)
    
    // Test connection
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    return &SupabaseClient{db: db}, nil
}

// Query executes a query that returns rows
func (sc *SupabaseClient) Query(query string, args ...interface{}) (*sql.Rows, error) {
    return sc.db.Query(query, args...)
}

// QueryRow executes a query that returns at most one row
func (sc *SupabaseClient) QueryRow(query string, args ...interface{}) *sql.Row {
    return sc.db.QueryRow(query, args...)
}

// Exec executes a query that doesn't return rows
func (sc *SupabaseClient) Exec(query string, args ...interface{}) (sql.Result, error) {
    return sc.db.Exec(query, args...)
}

// BeginTx starts a transaction
func (sc *SupabaseClient) BeginTx() (*sql.Tx, error) {
    return sc.db.Begin()
}

// Close closes the database connection
func (sc *SupabaseClient) Close() error {
    return sc.db.Close()
}

// Health checks if the database connection is healthy
func (sc *SupabaseClient) Health() error {
    return sc.db.Ping()
}

// DB exposes the underlying *sql.DB so services can use it directly.
func (sc *SupabaseClient) DB() *sql.DB {
    return sc.db
}
