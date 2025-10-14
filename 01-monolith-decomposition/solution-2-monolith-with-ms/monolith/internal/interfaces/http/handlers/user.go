package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"monolith/internal/usecase"
)

type UserHandler struct {
	userUseCase *usecase.UserUseCase
}

func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

type registerRequest struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Birthdate  string `json:"birthdate"`
	Biography  string `json:"biography"`
	City       string `json:"city"`
	Password   string `json:"password"`
}

type registerResponse struct {
	UserID string `json:"user_id"`
}

func (h *UserHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	var req registerRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Convert to use case request
	useCaseReq := &usecase.RegisterUserRequest{
		FirstName:  req.FirstName,
		SecondName: req.SecondName,
		Birthdate:  req.Birthdate,
		Biography:  req.Biography,
		City:       req.City,
		Password:   req.Password,
	}

	// Call use case
	response, err := h.userUseCase.RegisterUser(useCaseReq)
	if err != nil {
		switch err {
		case usecase.ErrInvalidData:
			http.Error(w, "invalid data", http.StatusBadRequest)
		case usecase.ErrInvalidBirthdate:
			http.Error(w, "invalid birthdate", http.StatusBadRequest)
		default:
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(registerResponse{UserID: response.UserID})
}
