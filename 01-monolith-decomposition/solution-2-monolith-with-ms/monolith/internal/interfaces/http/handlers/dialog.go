package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"monolith/internal/client"
)

type DialogHandler struct {
	dialogClient *client.DialogClient
}

func NewDialogHandler(dialogClient *client.DialogClient) *DialogHandler {
	return &DialogHandler{dialogClient: dialogClient}
}

type sendMessageRequest struct {
	Text string `json:"text"`
}

// HandleSendMessage handles POST /dialog/{user_id}/send
func (h *DialogHandler) HandleSendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get target user ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	targetUserID := pathParts[2] // /dialog/{user_id}/send

	// Parse request body
	var req sendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Get auth token from header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Call dialog service
	err := h.dialogClient.SendMessage(userID, targetUserID, req.Text, token)
	if err != nil {
		http.Error(w, "dialog service error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandleListMessages handles GET /dialog/{user_id}/list
func (h *DialogHandler) HandleListMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get target user ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	targetUserID := pathParts[2] // /dialog/{user_id}/list

	// Get auth token from header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Call dialog service
	messages, err := h.dialogClient.GetMessages(userID, targetUserID, token)
	if err != nil {
		http.Error(w, "dialog service error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
