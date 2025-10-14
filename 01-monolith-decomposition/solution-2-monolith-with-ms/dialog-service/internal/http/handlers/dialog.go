package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"dialog-service/internal/usecase"
)

type DialogHandler struct {
	dialogUseCase *usecase.DialogUseCase
}

func NewDialogHandler(dialogUseCase *usecase.DialogUseCase) *DialogHandler {
	return &DialogHandler{dialogUseCase: dialogUseCase}
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

	// Convert to use case request
	useCaseReq := &usecase.SendMessageRequest{
		FromUserID: userID,
		ToUserID:   targetUserID,
		Text:       req.Text,
	}

	// Call use case
	err := h.dialogUseCase.SendMessage(useCaseReq)
	if err != nil {
		switch err {
		case usecase.ErrInvalidData:
			http.Error(w, "text is required", http.StatusBadRequest)
		default:
			http.Error(w, "server error", http.StatusInternalServerError)
		}
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

	// Call use case
	messages, err := h.dialogUseCase.GetMessagesBetweenUsers(userID, targetUserID)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
