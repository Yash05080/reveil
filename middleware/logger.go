package middleware

import (
    "net/http"
    "time"
    
    "reveil-api/utils"
)

// Logger logs HTTP requests
func Logger(next http.Handler, logger *utils.Logger) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Create a response writer wrapper to capture status code
        wrapper := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        
        // Process request
        next.ServeHTTP(wrapper, r)
        
        // Log request details
        duration := time.Since(start)
        logger.Info("HTTP Request",
            "method", r.Method,
            "path", r.URL.Path,
            "status", wrapper.statusCode,
            "duration", duration.String(),
            "user_agent", r.UserAgent(),
            "remote_addr", r.RemoteAddr,
        )
    })
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
