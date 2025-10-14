package router

import (
	"net/http"
	"strings"

	"dialog-service/internal/http/handlers"
	"dialog-service/internal/http/middleware"
	"dialog-service/internal/usecase"
)

type Router struct {
	mux           *http.ServeMux
	dialogHandler *handlers.DialogHandler
	authUseCase   *usecase.AuthUseCase
}

func NewRouter(dialogHandler *handlers.DialogHandler, authUseCase *usecase.AuthUseCase) *Router {
	r := &Router{
		mux:           http.NewServeMux(),
		dialogHandler: dialogHandler,
		authUseCase:   authUseCase,
	}
	r.registerRoutes()
	return r
}

func (r *Router) registerRoutes() {
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
