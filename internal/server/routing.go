package server

import (
	"database/sql"
	"github.com/HDudz/SWIFT-Parser/internal/handler"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r *chi.Mux, db *sql.DB) {
	r.Get("/v1/swift-codes/{swift-code}", handler.GetCodeHandler(db))
	r.Get("/v1/swift-codes/country/{countryISO2code}", handler.GetCountryHandler(db))
	r.Post("/v1/swift-codes", handler.PostCodeHandler(db))
}
