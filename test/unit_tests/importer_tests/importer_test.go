package importer_tests

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/HDudz/SWIFT-Parser/internal/database"
	"github.com/HDudz/SWIFT-Parser/internal/models"
	"os"
	"regexp"
	"testing"
)

var db *sql.DB
var mock sqlmock.Sqlmock
var err error
var tmpFile *os.File

func TestMain(m *testing.M) {

	db, mock, err = sqlmock.New()
	if err != nil {
		panic("Error when creating sqlmock: " + err.Error())
	}
	defer db.Close()

	tmpFile, err = os.CreateTemp(os.TempDir(), "test_swift.csv")
	if err != nil {
		panic("Failed to create temp file: %v " + err.Error())
	}
	defer os.Remove(tmpFile.Name())

	testData := `COUNTRY ISO2 CODE,SWIFT CODE,CODE TYPE,NAME,ADDRESS,TOWN NAME,COUNTRY NAME,TIME ZONE
PL,ABCDEFGH,BIC11,Test Bank,"Test Address",Warsaw,POLAND,Europe/Warsaw`

	if _, err := tmpFile.WriteString(testData); err != nil {
		panic("Failed to save temp file: %v " + err.Error())
	}
	tmpFile.Close()

	fmt.Print("\n========Starting importers tests========\n\n")
	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestValidateImport(t *testing.T) {

	validRecord := models.MainSwift{
		Address:     "Valid Address",
		BankName:    "Valid Bank",
		CountryISO2: "US",
		CountryName: "United States",
		IsHQ:        true,
		Code:        "ABCDEFGH123",
	}

	wrongCode := validRecord
	wrongCode.Code = "123"

	wrongCountryISO2 := validRecord
	wrongCountryISO2.CountryISO2 = "U5"

	emptyCountryISO2 := validRecord
	emptyCountryISO2.CountryISO2 = ""

	wrongCountry := validRecord
	wrongCountry.CountryName = "   "

	emptyAddress := validRecord
	emptyAddress.Address = ""

	emptyBankName := validRecord
	emptyBankName.BankName = ""

	err := database.ValidateImport(validRecord)

	if err != nil {
		t.Errorf("ValidateImport shouldn't return error, returned: %v", err)
	} else {
		t.Logf("ValidateImport works correctly with valid record")
	}

	err = database.ValidateImport(wrongCode)
	if err == nil {
		t.Errorf("ValidateImport should return: invalid SWIFT code error, returned: %v", err)
	} else {
		t.Logf("ValidateImport works correctly with invalid code and returns: %v", err)
	}

	err = database.ValidateImport(wrongCountryISO2)
	if err == nil {
		t.Errorf("ValidateImport should return: invalid ISO2 error, returned: %v", err)
	} else {
		t.Logf("ValidateImport works correctly with invalid ISO2 and returns: %v", err)
	}

	err = database.ValidateImport(emptyCountryISO2)
	if err == nil {
		t.Errorf("ValidateImport should return: empty ISO2 error, returned: %v", err)
	} else {
		t.Logf("ValidateImport works correctly with empty ISO2 and returns: %v", err)
	}

	err = database.ValidateImport(wrongCountry)
	if err == nil {
		t.Errorf("ValidateImport should return: empty country name error, returned: %v", err)
	} else {
		t.Logf("ValidateImport works correctly with invalid contry name and returns: %v", err)
	}

	err = database.ValidateImport(emptyAddress)
	if err != nil {
		t.Errorf("ValidateImport shouldn't return error because address is not mandatory, returned: %v", err)
	} else {
		t.Logf("ValidateImport works correctly with empty address doesn't return error. Address is not mandatory.")
	}

	err = database.ValidateImport(emptyBankName)
	if err == nil {
		t.Errorf("ValidateImport should return: empty bank name error, returned: %v", err)
	} else {
		t.Logf("ValidateImport works correctly with empty bank name and returns: %v", err)
	}

}

func TestImportDataIfNeeded(t *testing.T) {
	rows1 := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM swiftTable`)).WillReturnRows(rows1)

	err, if_import := database.ImportDataIfNeeded(db, tmpFile.Name())

	if err != nil {
		t.Errorf("ImportDataIfNeeded returned error: %v", err)
	}

	if if_import == 1 {
		t.Errorf("ImportDataIfNeeded doesn't work correctly, it imported data when database wasn't empty ")
	} else {
		t.Logf("ImportDataIfNeeded works correctly, it didn't import data because database wasn't empty")
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
		t.Errorf("ImportDataIfNeeded doesn't work correctly; it didn't import when database was empty.")
	} else {
		t.Logf("ImportDataIfNeeded works correctly, it imported data because database was empty")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("SQL mock expectations weren't met: %v", err)
	}
}
