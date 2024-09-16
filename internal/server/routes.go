package server

import (
	"net/http"

	"zadanie-6105/internal/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router, tenderHandler *handlers.TenderHandler, bidHandler *handlers.BidHandler) {
	apiRouter := router.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}).Methods("GET")

	tenderHandler.RegisterRoutes(apiRouter)

	bidHandler.RegisterRoutes(apiRouter)
}
