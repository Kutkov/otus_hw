package handlers

import (
	"encoding/json"
	"net/http"

	"monolith/internal/usecase"
)

type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase}
}

type loginRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Convert to use case request
	useCaseReq := &usecase.LoginRequest{
		ID:       req.ID,
		Password: req.Password,
	}

	// Call use case
	response, err := h.authUseCase.Login(useCaseReq)
	if err != nil {
		switch err {
		case usecase.ErrInvalidData:
			http.Error(w, "invalid data", http.StatusBadRequest)
		case usecase.ErrUserNotFound, usecase.ErrInvalidPassword:
			http.Error(w, "not found", http.StatusNotFound)
		default:
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(loginResponse{Token: response.Token})
}
