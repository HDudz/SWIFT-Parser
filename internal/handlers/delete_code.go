package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"regexp"
)

func DeleteCodeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		code := chi.URLParam(r, "swift-code")

		if len(code) != 11 {
			http.Error(w, "Invalid SWIFT code format. Must be 11 characters long.", http.StatusBadRequest)
			return
		}

		if !regexp.MustCompile(`^[A-Z0-9]{8,11}$`).MatchString(code) {
			http.Error(w, "Invalid SWIFT code format. Must contain only uppercase letters and digits.", http.StatusBadRequest)
			return
		}

		var howManyFound int

		err := db.QueryRow(`SELECT COUNT(*) FROM swiftTable WHERE code = ?`, code).Scan(&howManyFound)

		if err != nil {
			log.Printf("Database error while checking SWIFT code: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if howManyFound == 0 {
			log.Printf("SWIFT code not found: %s", code)
			http.Error(w, "SWIFT code not found", http.StatusNotFound)
			return
		}

		_, err = db.Exec(`DELETE FROM swiftTable WHERE code = ?`, code)

		if err != nil {
			log.Printf("Database error while deleting SWIFT code %s: %v", code, err)
			http.Error(w, "Failed to delete data", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Swift code deleted successfully"})

	}

}
