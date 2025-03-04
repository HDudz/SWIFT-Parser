package handler

import (
	"database/sql"
	"encoding/json"
	"github.com/HDudz/SWIFT-Parser/internal/models"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func GetCountryHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ISO := chi.URLParam(r, "countryISO2code")

		var country models.CountryModel
		err := db.QueryRow(`SELECT country_iso2, country_name FROM swiftTable WHERE country_iso2 = ?`, ISO).
			Scan(&country.CountryISO2, &country.CountryName)

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
								   FROM swiftTable WHERE country_iso2 = ?`, country.CountryISO2)
		if err != nil {
			log.Println("DB error:", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		swiftCodes := []models.SubSwift{}

		for rows.Next() {
			var swiftCode models.SubSwift
			err = rows.Scan(&swiftCode.CountryISO2, &swiftCode.Code, &swiftCode.BankName, &swiftCode.Address, &swiftCode.IsHQ)
			if err != nil {
				log.Println("DB error:", err)
				break
			}
			swiftCodes = append(swiftCodes, swiftCode)
		}

		country.SwiftCodes = &swiftCodes
		response := country

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
