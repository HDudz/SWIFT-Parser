package services

import (
	"database/sql"
	"github.com/HDudz/SWIFT-Parser/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func LoadRoutes(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	swiftSubRouter := chi.NewRouter()

	swiftSubRouter.Get("/{swift-code}", handlers.GetCodeHandler(db))
	swiftSubRouter.Get("/country/{countryISO2code}", handlers.GetCountryHandler(db))
	swiftSubRouter.Post("/", handlers.PostCodeHandler(db))
	swiftSubRouter.Delete("/{swift-code}", handlers.DeleteCodeHandler(db))

	r.Mount("/v1/swift-codes", swiftSubRouter)

	return r
}
