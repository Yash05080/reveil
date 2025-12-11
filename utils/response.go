package utils

import (
    "encoding/json"
    "net/http"
    "time"
)

// StandardResponse represents the standard API response format
type StandardResponse struct {
    Success    bool        `json:"success"`
    StatusCode int         `json:"status_code"`
    Data       interface{} `json:"data,omitempty"`
    Error      string      `json:"error,omitempty"`
    ErrorCode  string      `json:"error_code,omitempty"`
    Details    interface{} `json:"details,omitempty"`
    Timestamp  time.Time   `json:"timestamp"`
}



// ErrorResponse writes an error JSON response
func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
    ErrorResponseWithCode(w, statusCode, message, "")
}



// ErrorResponseWithDetails writes an error JSON response with additional details
func ErrorResponseWithDetails(w http.ResponseWriter, statusCode int, message, errorCode string, details interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    response := StandardResponse{
        Success:    false,
        StatusCode: statusCode,
        Error:      message,
        ErrorCode:  errorCode,
        Details:    details,
        Timestamp:  time.Now(),
    }
    
    json.NewEncoder(w).Encode(response)
}

// JSONResponse writes a generic JSON response
func JSONResponse(w http.ResponseWriter, statusCode int, body interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(body)
}



// Standard API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

type APIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]string      `json:"details,omitempty"`
}

// SuccessResponse writes a success JSON response
func SuccessResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    data,
	})
}

// ErrorResponseWithCode writes an error with a specific code
func ErrorResponseWithCode(w http.ResponseWriter, status int, message string, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}

// ValidationErrorResponse formats validator errors
func ValidationErrorResponse(w http.ResponseWriter, err error) {
	errors := ParseValidationErrors(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Error: &APIError{
			Code:    "validation_error",
			Message: "Validation failed",
			Details: errors,
		},
	})
}

