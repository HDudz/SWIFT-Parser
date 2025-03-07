package handler

import (
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func DeleteCodeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		code := chi.URLParam(r, "swift-code")

		_, err := db.Exec(`DELETE FROM swiftTable WHERE code = ?`, code)

		if err != nil {
			log.Printf("Database error: %v", err)
			http.Error(w, "Failed to delete data", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Swift code deleted successfully"})

	}

}
