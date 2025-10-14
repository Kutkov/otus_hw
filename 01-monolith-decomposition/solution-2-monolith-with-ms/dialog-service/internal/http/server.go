package http

import (
	"database/sql"
	"net/http"

	"dialog-service/internal/http/handlers"
	"dialog-service/internal/http/router"
	"dialog-service/internal/repository"
	"dialog-service/internal/usecase"
)

type Server struct {
	db            *sql.DB
	router        *router.Router
	dialogHandler *handlers.DialogHandler
}

func NewServer(database *sql.DB) *Server {
	// Initialize repositories
	dialogRepo := repository.NewDialogRepository(database)

	// Initialize use cases
	dialogUseCase := usecase.NewDialogUseCase(dialogRepo)

	// Create a mock auth repository that will be replaced by HTTP calls to monolith
	mockAuthRepo := &MockAuthRepository{}
	authUseCase := usecase.NewAuthUseCase(mockAuthRepo)

	// Initialize handlers
	dialogHandler := handlers.NewDialogHandler(dialogUseCase)

	s := &Server{
		db:            database,
		dialogHandler: dialogHandler,
		router:        router.NewRouter(dialogHandler, authUseCase),
	}
	return s
}

func (s *Server) Handler() http.Handler {
	return s.router.Handler()
}

// MockAuthRepository is a placeholder that will be replaced by HTTP calls to monolith
type MockAuthRepository struct{}

func (m *MockAuthRepository) GetUserIDByToken(token string) (string, error) {
	// This will be replaced by actual HTTP call to monolith
	// For now, just return a mock user ID
	return "mock-user-id", nil
}
