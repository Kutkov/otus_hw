package usecase

import (
	"time"

	"monolith/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	authRepo *repository.AuthRepository
}

func NewAuthUseCase(authRepo *repository.AuthRepository) *AuthUseCase {
	return &AuthUseCase{authRepo: authRepo}
}

type LoginRequest struct {
	ID       string
	Password string
}

type LoginResponse struct {
	Token string
}

func (uc *AuthUseCase) Login(req *LoginRequest) (*LoginResponse, error) {
	// Validate required fields
	if req.ID == "" || req.Password == "" {
		return nil, ErrInvalidData
	}

	// Get password hash
	passwordHash, err := uc.authRepo.GetPasswordHash(req.ID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword(passwordHash, []byte(req.Password)); err != nil {
		return nil, ErrInvalidPassword
	}

	// Generate token
	token := uuid.New().String()
	createdAt := time.Now().UTC().Format("2006-01-02T15:04:05Z07:00")

	// Store token
	if err := uc.authRepo.CreateToken(token, req.ID, createdAt); err != nil {
		return nil, err
	}

	return &LoginResponse{Token: token}, nil
}

func (uc *AuthUseCase) ValidateToken(token string) (string, error) {
	userID, err := uc.authRepo.GetUserIDByToken(token)
	if err != nil {
		return "", ErrUnauthorized
	}
	return userID, nil
}
