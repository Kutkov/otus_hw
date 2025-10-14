package router

import (
	"net/http"
	"strings"

	"monolith/internal/interfaces/http/handlers"
	"monolith/internal/interfaces/http/middleware"
	"monolith/internal/usecase"
)

type Router struct {
	mux           *http.ServeMux
	authHandler   *handlers.AuthHandler
	userHandler   *handlers.UserHandler
	dialogHandler *handlers.DialogHandler
	authUseCase   *usecase.AuthUseCase
}

func NewRouter(authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, dialogHandler *handlers.DialogHandler, authUseCase *usecase.AuthUseCase) *Router {
	r := &Router{
		mux:           http.NewServeMux(),
		authHandler:   authHandler,
		userHandler:   userHandler,
		dialogHandler: dialogHandler,
		authUseCase:   authUseCase,
	}
	r.registerRoutes()
	return r
}

func (r *Router) registerRoutes() {
	r.mux.HandleFunc("/login", r.authHandler.HandleLogin)
	r.mux.HandleFunc("/user/register", r.userHandler.HandleRegister)

	// Dialog routes with authentication middleware
	authMiddleware := middleware.AuthMiddleware(r.authUseCase)

	// Create dialog handlers with auth middleware
	sendHandler := authMiddleware(http.HandlerFunc(r.dialogHandler.HandleSendMessage))
	listHandler := authMiddleware(http.HandlerFunc(r.dialogHandler.HandleListMessages))

	// Register dialog routes
	r.mux.HandleFunc("/dialog/", func(w http.ResponseWriter, req *http.Request) {
		if strings.HasSuffix(req.URL.Path, "/send") {
			sendHandler.ServeHTTP(w, req)
		} else if strings.HasSuffix(req.URL.Path, "/list") {
			listHandler.ServeHTTP(w, req)
		} else {
			http.NotFound(w, req)
		}
	})
}

func (r *Router) Handler() http.Handler {
	return r.mux
}
