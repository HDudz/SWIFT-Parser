package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/HDudz/SWIFT-Parser/internal/models"
	"log"
	"net/http"
	"strings"
)

func PostCodeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var rawData map[string]json.RawMessage

		if err := json.NewDecoder(r.Body).Decode(&rawData); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		requiredFields := []string{"address", "bankName", "countryISO2", "countryName", "isHeadquarter", "swiftCode"}

		for _, field := range requiredFields {
			_, exists := rawData[field]
			if !exists {
				http.Error(w, fmt.Sprintf("Missing required field: %s", field), http.StatusBadRequest)
				return
			}
		}

		var newCode models.MainSwift

		if err := json.Unmarshal(rawData["address"], &newCode.Address); err != nil || newCode.Address == "" {
			http.Error(w, "Invalid or missing 'address'", http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(rawData["bankName"], &newCode.BankName); err != nil || newCode.BankName == "" {
			http.Error(w, "Invalid or missing 'bankName'", http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(rawData["countryISO2"], &newCode.CountryISO2); err != nil || newCode.CountryISO2 == "" {
			http.Error(w, "Invalid or missing 'countryISO2'", http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(rawData["countryName"], &newCode.CountryName); err != nil || newCode.CountryName == "" {
			http.Error(w, "Invalid or missing 'countryName'", http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(rawData["swiftCode"], &newCode.Code); err != nil || newCode.Code == "" {
			http.Error(w, "Invalid or missing 'swiftCode'", http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(rawData["isHeadquarter"], &newCode.IsHQ); err != nil {
			http.Error(w, "Invalid 'isHeadquarter' field, must be true or false", http.StatusBadRequest)
			return
		}

		_, err := db.Exec(`INSERT INTO swiftTable (country_iso2, code, bank_name, address, country_name, is_hq) 
			VALUES (?, ?, ?, ?, ?, ?)`, strings.ToUpper(newCode.CountryISO2), newCode.Code, newCode.BankName, newCode.Address, strings.ToUpper(newCode.CountryName), newCode.IsHQ)

		if err != nil {
			log.Printf("Database insert error: %v", err)
			http.Error(w, "Failed to insert data", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Swift code inserted successfully"})

	}

}
