package handlers

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
		var howManyFound int

		err := db.QueryRow(`SELECT COUNT(*) FROM swiftTable WHERE code = ?`, code).Scan(&howManyFound)

		if err != nil {
			log.Printf("Database error: %v", err)
			http.Error(w, "Failed to find code", http.StatusNotFound)
			return
		}

		if howManyFound < 1 {
			log.Printf("Code not found")
			http.Error(w, "Code not found", http.StatusNotFound)
			return
		}

		_, err = db.Exec(`DELETE FROM swiftTable WHERE code = ?`, code)

		if err != nil {
			log.Printf("Database error: %v", err)
			http.Error(w, "Failed to delete data", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Swift code deleted successfully"})

	}

}
