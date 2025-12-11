package config

import "time"

const (
    // API Endpoints
    APIPrefix           = "/api"
    HealthEndpoint      = "/health"
    PostsEndpoint       = "/communities/{communityId}/posts"
    PostsStreamEndpoint = "/communities/{communityId}/posts/stream"
    
    // Content Limits
    MaxContentLength = 5000
    MinContentLength = 1
    
    // Moderation
    LightModerationTimeout = 500 * time.Millisecond
    HeavyModerationTimeout = 30 * time.Second
    
    // SSE Configuration
    SSETimeout = 30 * time.Minute
    DefaultLimit = 10
    
    // Rate Limiting
    PostsPerHour = 10
    MaxSSEConnections = 5
    
    // Database
    MaxDBConnections = 20
    
    // Error Codes
    ErrorValidation     = "VALIDATION_ERROR"
    ErrorAuthentication = "AUTHENTICATION_ERROR" 
    ErrorAuthorization  = "AUTHORIZATION_ERROR"
    ErrorContentBlocked = "CONTENT_BLOCKED"
    ErrorNotFound       = "NOT_FOUND"
    ErrorInternal       = "INTERNAL_ERROR"
    ErrorRateLimit      = "RATE_LIMIT_EXCEEDED"
)
