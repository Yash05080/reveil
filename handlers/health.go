package handlers

import (
    //"encoding/json"
    "net/http"
    "time"
    
    "reveil-api/db"
    "reveil-api/utils"
)

// HealthHandler handles health check requests
type HealthHandler struct {
    db *db.SupabaseClient
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(database *db.SupabaseClient) *HealthHandler {
    return &HealthHandler{
        db: database,
    }
}

// HealthResponse represents the health check response
type HealthResponse struct {
    Status    string    `json:"status"`
    Timestamp time.Time `json:"timestamp"`
    Version   string    `json:"version"`
    Database  string    `json:"database"`
    Uptime    string    `json:"uptime"`
}

var startTime = time.Now()

// HealthCheck handles GET /health
func (hh *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
    // Check database connectivity
    dbStatus := "healthy"
    if err := hh.db.Health(); err != nil {
        dbStatus = "unhealthy"
        
        // If database is down, return 503
        response := HealthResponse{
            Status:    "unhealthy",
            Timestamp: time.Now(),
            Version:   "1.0.0",
            Database:  dbStatus,
            Uptime:    time.Since(startTime).String(),
        }
        
        utils.JSONResponse(w, http.StatusServiceUnavailable, response)
        return
    }
    
    // All systems healthy
    response := HealthResponse{
        Status:    "healthy",
        Timestamp: time.Now(),
        Version:   "1.0.0",
        Database:  dbStatus,
        Uptime:    time.Since(startTime).String(),
    }
    
    utils.JSONResponse(w, http.StatusOK, response)
}



/*

curl http://localhost:8080/health | jq '.database'
curl http://localhost:8080/nonexistent    
curl -i -X OPTIONS http://localhost:8080/health   
curl http://localhost:8080/health              

*/