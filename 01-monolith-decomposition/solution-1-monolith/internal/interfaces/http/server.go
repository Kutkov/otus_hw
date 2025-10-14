package http

import (
	"database/sql"
	"net/http"

	"monolith/internal/interfaces/http/handlers"
	"monolith/internal/interfaces/http/router"
	"monolith/internal/repository"
	"monolith/internal/usecase"
)

type Server struct {
	db            *sql.DB
	router        *router.Router
	authHandler   *handlers.AuthHandler
	userHandler   *handlers.UserHandler
	dialogHandler *handlers.DialogHandler
}

func NewServer(database *sql.DB) *Server {
	// Initialize repositories
	userRepo := repository.NewUserRepository(database)
	authRepo := repository.NewAuthRepository(database)
	dialogRepo := repository.NewDialogRepository(database)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo)
	authUseCase := usecase.NewAuthUseCase(authRepo)
	dialogUseCase := usecase.NewDialogUseCase(dialogRepo, userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authUseCase)
	userHandler := handlers.NewUserHandler(userUseCase)
	dialogHandler := handlers.NewDialogHandler(dialogUseCase)

	s := &Server{
		db:            database,
		authHandler:   authHandler,
		userHandler:   userHandler,
		dialogHandler: dialogHandler,
		router:        router.NewRouter(authHandler, userHandler, dialogHandler, authUseCase),
	}
	return s
}

func (s *Server) Handler() http.Handler {
	return s.router.Handler()
}
