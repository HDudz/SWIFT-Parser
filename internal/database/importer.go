package database

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/HDudz/SWIFT-Parser/internal/models"
	"os"
	"regexp"
	"strings"
)

func ValidateImport(record models.MainSwift) error {

	var err error
	if len(record.CountryISO2) != 2 {
		err = errors.Join(err, fmt.Errorf("iso2 code length must be 2, given ISO2: \"%s\"", record.CountryISO2))
	}
	isOk, _ := regexp.MatchString(`^[A-Z]*$`, record.CountryISO2)
	if !isOk {
		err = errors.Join(err, fmt.Errorf("iso2 code can only contain letters, given ISO2: \"%s\"", record.CountryISO2))
	}

	if len(record.Code) < 8 || len(record.Code) > 11 {
		err = errors.Join(err, fmt.Errorf("SWIFT code length must be 8-11, given code: \"%s\"", record.Code))
	}

	if strings.TrimSpace(record.CountryName) == "" {
		err = errors.Join(err, fmt.Errorf("country name must not be empty"))
	}

	if strings.TrimSpace(record.BankName) == "" {
		err = errors.Join(err, fmt.Errorf("bank name must not be empty"))
	}

	return err
}

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

	var newCode models.MainSwift

	for i, record := range records[1:] {
		recordNumber := i + 1
		newCode.CountryISO2 = strings.ToUpper(record[0])
		newCode.Code = record[1]
		newCode.BankName = record[3]
		newCode.Address = record[4]
		newCode.CountryName = strings.ToUpper(record[6])
		newCode.IsHQ = strings.HasSuffix(newCode.Code, "XXX")
		err := ValidateImport(newCode)
		if err == nil {
			_, err := db.Exec(`
            INSERT INTO swiftTable (country_iso2, code, bank_name, address, country_name, is_hq) 
            VALUES (?, ?, ?, ?, ?, ?)`,
				newCode.CountryISO2, newCode.Code, newCode.BankName, newCode.Address, newCode.CountryName, newCode.IsHQ)

			if err != nil {
				return fmt.Errorf("couldn't insert data into db, record: %s: error: %v\n", newCode.Code, err)
			}
		} else {
			fmt.Printf("Skipped record number %d, error: %v\n", recordNumber, err)
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
