package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"reveil-api/config"
	"reveil-api/models"
	"reveil-api/services"
	"reveil-api/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type PostHandler struct {
	postService *services.PostService
	sse         *services.SSEService
}

// NewPostHandler creates a new PostHandler
func NewPostHandler(ps *services.PostService, sse *services.SSEService) *PostHandler {
	return &PostHandler{
		postService: ps,
		sse:         sse,
	}
}

// RegisterPostRoutes attaches post routes to router
func (h *PostHandler) RegisterPostRoutes(r *mux.Router, validator *utils.Validator) {
	// The router 'r' passed here is already a subrouter on "/api"
	// So we register "/communities/{community_id}/posts" directly
	r.HandleFunc("/communities/{community_id}/posts", h.createPost(validator)).Methods(http.MethodPost)
	r.HandleFunc("/communities/{community_id}/posts", h.listPosts(validator)).Methods(http.MethodGet)
	r.HandleFunc("/communities/{community_id}/posts/stream", h.streamPosts()).Methods(http.MethodGet)
	r.HandleFunc("/posts/{post_id}", h.updatePost(validator)).Methods(http.MethodPut)
	r.HandleFunc("/posts/{post_id}", h.deletePost()).Methods(http.MethodDelete)
}

// streamPosts handles SSE connection for real-time posts
// GET /api/communities/{community_id}/posts/stream
func (h *PostHandler) streamPosts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		communityIDStr := vars["community_id"]
		communityID, err := uuid.Parse(communityIDStr)
		if err != nil {
			utils.ErrorResponseWithCode(w, http.StatusBadRequest, "Invalid community id", config.ErrorValidation)
			return
		}

		// Set SSE headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*") // Adjust in production

		clientChan := h.sse.Subscribe(communityID)
		defer h.sse.Unsubscribe(communityID, clientChan)

		// Notify client of connection
		fmt.Fprintf(w, "data: {\"status\": \"connected\", \"community_id\": \"%s\"}\n\n", communityID)

		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}

		// Listen for events or context cancellation
		for {
			select {
			case event := <-clientChan:
				data, _ := json.Marshal(event)
				fmt.Fprintf(w, "data: %s\n\n", data)
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			case <-r.Context().Done():
				return
			}
		}
	}
}

// createPost handles POST /api/communities/{community_id}/posts
func (h *PostHandler) createPost(validator *utils.Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// user_id and community_id from JWT/context
		userIDVal := r.Context().Value("user_id")
		if userIDVal == nil {
			utils.ErrorResponseWithCode(w, http.StatusUnauthorized, "Unauthorized", config.ErrorAuthentication)
			return
		}
		userIDStr, ok := userIDVal.(string)
		if !ok {
			utils.ErrorResponseWithCode(w, http.StatusUnauthorized, "Invalid user context", config.ErrorAuthentication)
			return
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			utils.ErrorResponseWithCode(w, http.StatusUnauthorized, "Invalid user id", config.ErrorAuthentication)
			return
		}

		vars := mux.Vars(r)
		communityIDStr := vars["community_id"]
		communityID, err := uuid.Parse(communityIDStr)
		if err != nil {
			utils.ErrorResponseWithCode(w, http.StatusBadRequest, "Invalid community id", config.ErrorValidation)
			return
		}

		var req models.CreatePostRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorResponseWithCode(w, http.StatusBadRequest, "Invalid JSON body", config.ErrorValidation)
			return
		}

		if err := validator.ValidateStruct(req); err != nil {
			utils.ValidationErrorResponse(w, err)
			return
		}

		post, err := h.postService.CreatePost(r.Context(), communityID, userID, req)
		if err != nil {
			utils.ErrorResponseWithCode(w, http.StatusInternalServerError, "Failed to create post", config.ErrorInternal)
			return
		}

		utils.SuccessResponse(w, http.StatusCreated, post)
	}
}

// listPosts handles GET /api/communities/{community_id}/posts
func (h *PostHandler) listPosts(validator *utils.Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		communityIDStr := vars["community_id"]
		communityID, err := uuid.Parse(communityIDStr)
		if err != nil {
			utils.ErrorResponseWithCode(w, http.StatusBadRequest, "Invalid community id", config.ErrorValidation)
			return
		}

		limitStr := r.URL.Query().Get("limit")
		beforeStr := r.URL.Query().Get("before")

		q := models.ListPostsQuery{}
		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil {
				q.Limit = l
			}
		}
		if beforeStr != "" {
			if t, err := time.Parse(time.RFC3339, beforeStr); err == nil {
				q.Before = &t
			}
		}

		userIDFilterStr := r.URL.Query().Get("user_id")
		if userIDFilterStr != "" {
			if uid, err := uuid.Parse(userIDFilterStr); err == nil {
				q.UserID = &uid
			}
		}

		contentTypeStr := r.URL.Query().Get("content_type")
		if contentTypeStr != "" {
			q.ContentType = &contentTypeStr
		}

		posts, err := h.postService.ListPosts(r.Context(), communityID, q)
		if err != nil {
			utils.ErrorResponseWithCode(w, http.StatusInternalServerError, "Failed to fetch posts", config.ErrorInternal)
			return
		}

		utils.SuccessResponse(w, http.StatusOK, posts)
	}
}

// updatePost handles PUT /api/posts/{post_id}
func (h *PostHandler) updatePost(validator *utils.Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDVal := r.Context().Value("user_id")
		if userIDVal == nil {
			utils.ErrorResponseWithCode(w, http.StatusUnauthorized, "Unauthorized", config.ErrorAuthentication)
			return
		}
		userID, _ := uuid.Parse(userIDVal.(string))

		vars := mux.Vars(r)
		postID, err := uuid.Parse(vars["post_id"])
		if err != nil {
			utils.ErrorResponseWithCode(w, http.StatusBadRequest, "Invalid post id", config.ErrorValidation)
			return
		}

		var req models.UpdatePostRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorResponseWithCode(w, http.StatusBadRequest, "Invalid JSON body", config.ErrorValidation)
			return
		}

		if err := validator.ValidateStruct(req); err != nil {
			utils.ValidationErrorResponse(w, err)
			return
		}

		post, err := h.postService.UpdatePost(r.Context(), postID, userID, req)
		if err != nil {
			if err.Error() == "unauthorized" {
				utils.ErrorResponseWithCode(w, http.StatusForbidden, "Not authorized to update this post", config.ErrorAuthorization)
				return
			}
			if err.Error() == "post not found" {
				utils.ErrorResponseWithCode(w, http.StatusNotFound, "Post not found", config.ErrorNotFound)
				return
			}
			utils.ErrorResponseWithCode(w, http.StatusInternalServerError, "Failed to update post", config.ErrorInternal)
			return
		}

		utils.SuccessResponse(w, http.StatusOK, post)
	}
}

// deletePost handles DELETE /api/posts/{post_id}
func (h *PostHandler) deletePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDVal := r.Context().Value("user_id")
		if userIDVal == nil {
			utils.ErrorResponseWithCode(w, http.StatusUnauthorized, "Unauthorized", config.ErrorAuthentication)
			return
		}
		userID, _ := uuid.Parse(userIDVal.(string))

		vars := mux.Vars(r)
		postID, err := uuid.Parse(vars["post_id"])
		if err != nil {
			utils.ErrorResponseWithCode(w, http.StatusBadRequest, "Invalid post id", config.ErrorValidation)
			return
		}

		err = h.postService.DeletePost(r.Context(), postID, userID)
		if err != nil {
			if err.Error() == "post not found or unauthorized" {
				// We don't distinguish to avoid leaking existence, or we can use 404/403
				utils.ErrorResponseWithCode(w, http.StatusNotFound, "Post not found", config.ErrorNotFound)
				return
			}
			utils.ErrorResponseWithCode(w, http.StatusInternalServerError, "Failed to delete post", config.ErrorInternal)
			return
		}

		utils.SuccessResponse(w, http.StatusOK, map[string]string{"status": "deleted"})
	}
}
