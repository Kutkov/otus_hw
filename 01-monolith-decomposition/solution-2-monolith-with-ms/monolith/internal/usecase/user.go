package usecase

import (
	"time"

	"monolith/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepo *repository.UserRepository
}

func NewUserUseCase(userRepo *repository.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

type RegisterUserRequest struct {
	FirstName  string
	SecondName string
	Birthdate  string
	Biography  string
	City       string
	Password   string
}

type RegisterUserResponse struct {
	UserID string
}

func (uc *UserUseCase) RegisterUser(req *RegisterUserRequest) (*RegisterUserResponse, error) {
	// Validate required fields
	if req.FirstName == "" || req.SecondName == "" || req.Birthdate == "" || req.Password == "" {
		return nil, ErrInvalidData
	}

	// Validate birthdate format
	if _, err := time.Parse("2006-01-02", req.Birthdate); err != nil {
		return nil, ErrInvalidBirthdate
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	userID := uuid.New().String()
	user := &repository.User{
		ID:           userID,
		FirstName:    req.FirstName,
		SecondName:   req.SecondName,
		Birthdate:    req.Birthdate,
		Biography:    req.Biography,
		City:         req.City,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC().Format("2006-01-02T15:04:05Z07:00"),
	}

	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &RegisterUserResponse{UserID: userID}, nil
}

func (uc *UserUseCase) UserExists(userID string) (bool, error) {
	return uc.userRepo.Exists(userID)
}
