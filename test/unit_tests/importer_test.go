package unit_tests

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/HDudz/SWIFT-Parser/internal/database"
	"os"
	"regexp"
	"testing"
)

func TestImportFromCSV(t *testing.T) {

	fmt.Print("\n========Starting Importer Unit tests========\n\n")

	db, mock, err := sqlmock.New()
	if err != nil {
		panic("Error when creating sqlmock: " + err.Error())
	}
	defer db.Close()

	tmpFile, err := os.CreateTemp(os.TempDir(), "test_swift.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testData := `COUNTRY ISO2 CODE,SWIFT CODE,CODE TYPE,NAME,ADDRESS,TOWN NAME,COUNTRY NAME,TIME ZONE
PL,ABCDEFGH,BIC11,Test Bank,"Test Address",Warsaw,POLAND,Europe/Warsaw`
	if _, err := tmpFile.WriteString(testData); err != nil {
		t.Fatalf("Failed to save temp file: %v", err)
	}
	tmpFile.Close()

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO swiftTable (country_iso2, code, bank_name, address, country_name, is_hq) 
            VALUES (?, ?, ?, ?, ?, ?)`)).
		WithArgs("PL", "ABCDEFGH", "Test Bank", "Test Address", "POLAND", false).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = database.ImportFromCSV(db, tmpFile.Name())
	if err != nil {
		t.Errorf("ImportFromCSV returned error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("SQL mock expectations weren't met: %v", err)
	}
}

func TestImportDataIfNeeded(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		panic("Error when creating sqlmock: " + err.Error())
	}
	defer db.Close()

	tmpFile, err := os.CreateTemp(os.TempDir(), "test_swift.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testData := `COUNTRY ISO2 CODE,SWIFT CODE,CODE TYPE,NAME,ADDRESS,TOWN NAME,COUNTRY NAME,TIME ZONE
PL,ABCDEFGH,BIC11,Test Bank,"Test Address",Warsaw,POLAND,Europe/Warsaw`

	if _, err := tmpFile.WriteString(testData); err != nil {
		t.Fatalf("Failed to save temp file: %v", err)
	}
	tmpFile.Close()

	rows1 := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM swiftTable`)).WillReturnRows(rows1)

	err, if_import := database.ImportDataIfNeeded(db, tmpFile.Name())

	if err != nil {
		t.Errorf("ImportDataIfNeeded returned error: %v", err)
	}

	if if_import != 0 {
		t.Errorf("Data import status error: expected: 0, received: %v ", if_import)
	}

	rows2 := sqlmock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM swiftTable`)).WillReturnRows(rows2)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO swiftTable (country_iso2, code, bank_name, address, country_name, is_hq) 
            VALUES (?, ?, ?, ?, ?, ?)`)).
		WithArgs("PL", "ABCDEFGH", "Test Bank", "Test Address", "POLAND", false).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err, if_import = database.ImportDataIfNeeded(db, tmpFile.Name())
	if err != nil {
		t.Errorf("ImportDataIfNeeded returned error: %v", err)
	}

	if if_import != 1 {
		t.Errorf("Data import status error: expected: 1, received: %v ", if_import)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("SQL mock expectations weren't met: %v", err)
	}
}
