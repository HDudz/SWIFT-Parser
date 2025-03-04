package handler

import (
	"database/sql"
	"encoding/json"
	"github.com/HDudz/SWIFT-Parser/internal/models"
	"log"
	"net/http"
	"strings"
)

func PostCodeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newCode models.MainSwift

		if err := json.NewDecoder(r.Body).Decode(&newCode); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		log.Printf(strings.ToUpper(newCode.CountryISO2), newCode.Code, newCode.BankName, newCode.Address, strings.ToUpper(newCode.CountryName), newCode.IsHQ)
		_, err := db.Exec(`INSERT INTO swiftTable (country_iso2, code, bank_name, address, country_name, is_hq) 
			VALUES (?, ?, ?, ?, ?, ?)`, strings.ToUpper(newCode.CountryISO2), newCode.Code, newCode.BankName, newCode.Address, strings.ToUpper(newCode.CountryName), newCode.IsHQ)

		if err != nil {
			log.Printf("Database insert error: %v", err)
			http.Error(w, "Failed to insert data", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Swift code inserted successfully"})

	}

}
