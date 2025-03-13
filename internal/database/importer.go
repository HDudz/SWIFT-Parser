package database

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func ImportFromCSV(db *sql.DB, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("couldn't open file, error: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("couldn't read CSV, error: %v", err)
	}

	for _, record := range records[1:] {

		countryISO2 := strings.ToUpper(record[0])
		swiftCode := record[1]
		bankName := record[3]
		address := record[4]
		countryName := strings.ToUpper(record[6])
		isHQ := strings.HasSuffix(swiftCode, "XXX")

		_, err := db.Exec(`
            INSERT INTO swiftTable (country_iso2, code, bank_name, address, country_name, is_hq) 
            VALUES (?, ?, ?, ?, ?, ?)`,
			countryISO2, swiftCode, bankName, address, countryName, isHQ)

		if err != nil {
			return fmt.Errorf("couldn't insert data into db, record: %s: error: %v\n", swiftCode, err)
		}
	}
	return nil
}

func ImportDataIfNeeded(db *sql.DB, path string) (error, int) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM swiftTable").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check database: %v", err), -1
	}
	if count == 0 {
		err = ImportFromCSV(db, path)
		if err != nil {
			return fmt.Errorf("failed to import data from CSV: %v", err), -1
		}
		fmt.Println("Data has been imported")
		return nil, 1
	} else {
		fmt.Println("Data already exists, skipping import")
		return nil, 0
	}

}
