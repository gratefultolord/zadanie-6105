package server

import (
	"context"
	"net/http"
	"zadanie-6105/internal/config"
	"zadanie-6105/internal/handlers"
	"zadanie-6105/internal/repositories"
	"zadanie-6105/internal/services"

	"log"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, db *gorm.DB) (*Server, error) {
	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api").Subrouter()
	registerRoutes(apiRouter, db)

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	return &Server{httpServer: srv}, nil
}

func (s *Server) Start() error {
	log.Printf("Starting server on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}

func registerRoutes(router *mux.Router, db *gorm.DB) {
	tenderRepo := repositories.NewTenderRepository(db)
	bidRepo := repositories.NewBidRepository(db)

	tenderService := services.NewTenderService(tenderRepo)
	bidService := services.NewBidService(bidRepo)

	tenderHandler := handlers.NewTenderHandler(tenderService)
	bidHandler := handlers.NewBidHandler(bidService)

	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}).Methods("GET")

	tenderHandler.RegisterRoutes(router)
	bidHandler.RegisterRoutes(router)
}
