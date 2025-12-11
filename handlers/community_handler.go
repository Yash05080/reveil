package handlers

import (
	"encoding/json"
	"net/http"

	"reveil-api/config"
	"reveil-api/models"
	"reveil-api/services"
	"reveil-api/utils"

	"github.com/gorilla/mux"
)

type CommunityHandler struct {
	communityService *services.CommunityService
}

func NewCommunityHandler(cs *services.CommunityService) *CommunityHandler {
	return &CommunityHandler{communityService: cs}
}

func (h *CommunityHandler) RegisterCommunityRoutes(r *mux.Router, validator *utils.Validator) {
	r.HandleFunc("/communities", h.ListCommunities).Methods(http.MethodGet)
	r.HandleFunc("/communities", h.CreateCommunity(validator)).Methods(http.MethodPost)
}

func (h *CommunityHandler) ListCommunities(w http.ResponseWriter, r *http.Request) {
	communities, err := h.communityService.ListCommunities(r.Context())
	if err != nil {
		utils.ErrorResponseWithCode(w, http.StatusInternalServerError, "Failed to list communities", config.ErrorInternal)
		return
	}
	utils.SuccessResponse(w, http.StatusOK, communities)
}

func (h *CommunityHandler) CreateCommunity(validator *utils.Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.CreateCommunityRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorResponseWithCode(w, http.StatusBadRequest, "Invalid JSON", config.ErrorValidation)
			return
		}
		if err := validator.ValidateStruct(req); err != nil {
			utils.ValidationErrorResponse(w, err)
			return
		}

		c, err := h.communityService.CreateCommunity(r.Context(), req.Name, req.Description)
		if err != nil {
			utils.ErrorResponseWithCode(w, http.StatusInternalServerError, "Failed to create", config.ErrorInternal)
			return
		}
		utils.SuccessResponse(w, http.StatusCreated, c)
	}
}
