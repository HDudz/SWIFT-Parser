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
		return fmt.Errorf("Couldn't open file, error: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("Couldn't read CSV, error: %v", err)
	}

	for _, record := range records[1:] {

		fmt.Println(record)
		countryISO2 := strings.ToUpper(record[0])
		swiftCode := record[1]
		codeType := record[2]
		bankName := record[3]
		address := record[4]
		town := record[5]
		countryName := strings.ToUpper(record[6])
		timeZone := record[7]
		isHQ := strings.HasSuffix(swiftCode, "XXX")

		_, err := db.Exec(`
            INSERT INTO swiftTable (country_iso2, code, code_type, bank_name, address, town, country_name, time_zone, is_hq) 
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			countryISO2, swiftCode, codeType, bankName, address, town, countryName, timeZone, isHQ)

		if err != nil {
			return fmt.Errorf("Couldn't insert data into db, record: %s: error: %v\n", swiftCode, err)
		}
	}
	fmt.Println("Data imported successfully!")
	return nil
}
