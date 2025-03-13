package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/HDudz/SWIFT-Parser/internal/models"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func GetCodeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "swift-code")

		isHQ := strings.HasSuffix(code, "XXX")

		if len(code) != 11 {
			http.Error(w, "Invalid SWIFT code format. Must be 11 characters long.", http.StatusBadRequest)
			return
		}

		if !regexp.MustCompile(`^[A-Z0-9]{8,11}$`).MatchString(code) {
			http.Error(w, "Invalid SWIFT code format. Must contain only uppercase letters and digits.", http.StatusBadRequest)
			return
		}

		if isHQ {
			var hq models.MainSwift
			err := db.QueryRow(`SELECT country_iso2, code, bank_name, address, country_name, is_hq
								FROM swiftTable WHERE code = ?`, code).
				Scan(&hq.CountryISO2, &hq.Code, &hq.BankName, &hq.Address, &hq.CountryName, &hq.IsHQ)

			if err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, "Requested SWIFT Code not found.", http.StatusNotFound)
				} else {
					log.Println("DB error:", err)
					http.Error(w, "Internal error: ", http.StatusInternalServerError)
				}
				return
			}

			rows, err := db.Query(`SELECT country_iso2, code, bank_name, address, is_hq 
								   FROM swiftTable WHERE code != ? AND LEFT(code, 8) = ?`, hq.Code, hq.Code[:8])

			if err != nil {
				log.Println("DB error:", err)
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			branches := []models.SubSwift{}

			for rows.Next() {
				var branch models.SubSwift
				err = rows.Scan(&branch.CountryISO2, &branch.Code, &branch.BankName, &branch.Address, &branch.IsHQ)
				if err != nil {
					log.Println("DB error:", err)
					break
				}
				branches = append(branches, branch)
			}

			hq.Branches = &branches
			response := hq

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			var branch models.MainSwift
			err := db.QueryRow(`SELECT country_iso2, code, bank_name, address, country_name, is_hq
								FROM swiftTable WHERE code = ?`, code).
				Scan(&branch.CountryISO2, &branch.Code, &branch.BankName, &branch.Address, &branch.CountryName, &branch.IsHQ)

			if err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, "SWIFT code not found", http.StatusNotFound)
				} else {
					log.Println("DB error:", err)
					http.Error(w, "Internal error", http.StatusInternalServerError)
				}
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(branch)
		}
	}
}
