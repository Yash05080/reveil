package handlers

import (
	"net/http"

	"reveil-api/services"
	"reveil-api/utils"

	"github.com/gorilla/mux"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(as *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: as}
}

func (h *AuthHandler) RegisterAuthRoutes(r *mux.Router, validator *utils.Validator) {
	// Basic placeholder for now since we rely on external auth (Supabase/Firebase) usually,
	// but code structure implied local auth service.
	// r.HandleFunc("/login", h.Login).Methods("POST")
}

// Login placeholder
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	utils.SuccessResponse(w, http.StatusOK, map[string]string{"message": "Login not implemented (Placeholder)"})
}
