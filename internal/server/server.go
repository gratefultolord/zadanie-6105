package server

import (
	"context"
	"net/http"
	"time"

	"zadanie-6105/internal/config"
	"zadanie-6105/internal/handlers"
	"zadanie-6105/internal/middlewares"
	"zadanie-6105/internal/repositories"
	"zadanie-6105/internal/services"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, db *gorm.DB) (*Server, error) {
	tenderRepo := repositories.NewTenderRepository(db)
	bidRepo := repositories.NewBidRepository(db)
	employeeRepo := repositories.NewEmployeeRepository(db)
	organizationRepo := repositories.NewOrganizationRepository(db)

	tenderService := services.NewTenderService(tenderRepo)
	bidService := services.NewBidService(bidRepo, employeeRepo, organizationRepo)

	tenderHandler := handlers.NewTenderHandler(tenderService)
	bidHandler := handlers.NewBidHandler(bidService)

	router := mux.NewRouter()

	authMiddleware := middlewares.AuthMiddleware(db)

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(authMiddleware)

	tenderHandler.RegisterRoutes(apiRouter)
	bidHandler.RegisterRoutes(apiRouter)

	srv := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return &Server{
		httpServer: srv,
	}, nil
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
