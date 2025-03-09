package test

import (
	"context"
	"encoding/json"
	"github.com/HDudz/SWIFT-Parser/internal/handlers"
	"github.com/HDudz/SWIFT-Parser/internal/models"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
)

func TestGetCodeHandler_HQ(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic("Error when creating sqlmock: " + err.Error())
	}
	defer db.Close()

	if err := mock.ExpectationsWereMet(); err != nil {
		panic("There are unmet expectations: " + err.Error())
	}
	code := "ABCDEFXXX"
	queryMain := `
		SELECT country_iso2, code, bank_name, address, country_name, is_hq
		FROM swiftTable WHERE code = ?`
	rowsMain := sqlmock.NewRows([]string{"country_iso2", "code", "bank_name", "address", "country_name", "is_hq"}).
		AddRow("US", code, "HQ Bank", "HQ Address", "United States", true)
	mock.ExpectQuery(regexp.QuoteMeta(queryMain)).
		WithArgs(code).
		WillReturnRows(rowsMain)

	queryBranches := `
		SELECT country_iso2, code, bank_name, address, is_hq 
		FROM swiftTable WHERE code != ? AND LEFT(code, 8) = ?`
	prefix := code[:8]
	rowsBranches := sqlmock.NewRows([]string{"country_iso2", "code", "bank_name", "address", "is_hq"}).
		AddRow("US", "ABCDEF01", "Branch Bank", "Branch Address", false)
	mock.ExpectQuery(regexp.QuoteMeta(queryBranches)).
		WithArgs(code, prefix).
		WillReturnRows(rowsBranches)

	req := httptest.NewRequest("GET", "/v1/swift-codes/"+code, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("swift-code", code)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	h := handlers.GetCodeHandler(db)
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status: received %v, expected %v", status, http.StatusOK)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Wrong Content-Type: received %v, expected application/json", ct)
	}

	var resp models.MainSwift
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("Error when decoding response: %v", err)
	}

	if resp.Code != code {
		t.Errorf("Wrong code: expected %v, received %v", code, resp.Code)
	}
}

func TestGetCodeHandler_NonHQ(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic("Error when creating sqlmock: " + err.Error())
	}
	defer db.Close()

	code := "ABCDEF01"

	query := `
		SELECT country_iso2, code, bank_name, address, country_name, is_hq
		FROM swiftTable WHERE code = ?`
	rows := sqlmock.NewRows([]string{"country_iso2", "code", "bank_name", "address", "country_name", "is_hq"}).
		AddRow("US", code, "Branch Bank", "Branch Address", "United States", false)
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(code).
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/v1/swift-codes/"+code, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("swift-code", code)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	h := handlers.GetCodeHandler(db)
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status: received %v, expected %v", status, http.StatusOK)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Wrong Content-Type: received %v, expected application/json", ct)
	}

	var resp models.MainSwift
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("Error when decoding response: %v", err)
	}

	if resp.Code != code {
		t.Errorf("Wrong code: expected %v, received %v", code, resp.Code)
	}

}

func TestGetCountryHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic("Error when creating sqlmock: " + err.Error())
	}
	defer db.Close()

	ISO := "PL"

	queryCountry := `SELECT country_iso2, country_name FROM swiftTable WHERE country_iso2 = ?`
	rowsCountry := sqlmock.NewRows([]string{"country_iso2", "country_name"}).
		AddRow("PL", "Poland")
	mock.ExpectQuery(regexp.QuoteMeta(queryCountry)).
		WithArgs(ISO).
		WillReturnRows(rowsCountry)

	queryBank := `SELECT country_iso2, code, bank_name, address, is_hq FROM swiftTable WHERE country_iso2 = ?`
	rowsBank := sqlmock.NewRows([]string{"country_iso2", "code", "bank_name", "address", "is_hq"}).
		AddRow("PL", "ABCDEFGH", "Polish Bank", "Polish address", false)
	mock.ExpectQuery(regexp.QuoteMeta(queryBank)).
		WithArgs(ISO).
		WillReturnRows(rowsBank)

	req := httptest.NewRequest("GET", "/v1/swift-codes/country/"+ISO, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("countryISO2code", ISO)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	h := handlers.GetCountryHandler(db)
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status: received %v, expected %v", status, http.StatusOK)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Wrong Content-Type: received %v, expected application/json", ct)
	}

	var resp models.CountryModel
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("Error when decoding response: %v", err)
	}

	if resp.CountryISO2 != ISO {
		t.Errorf("Wrong ISO code: expected %v, received %v", ISO, resp.CountryISO2)
	}

	if resp.CountryName != "Poland" {
		t.Errorf("Wrong country name: expected 'Poland', received %v", resp.CountryName)
	}

	if len(*resp.SwiftCodes) != 1 {
		t.Errorf("Expected 1 swift code, received %d", len(*resp.SwiftCodes))
	}

	swiftCode := (*resp.SwiftCodes)[0]
	if swiftCode.Code != "ABCDEFGH" {
		t.Errorf("Wring swift code: received %v, expected ABCDEFGH", swiftCode.Code)
	}

}

func TestPostCodeHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic("Error when creating sqlmock: " + err.Error())
	}
	defer db.Close()

	reqBody := `{
		"countryISO2": "PL",
		"swiftCode": "ABCDEFGH",
		"bankName": "Polish Bank",
		"address": "Polish address",
		"countryName": "Poland",
		"isHeadquarter": true
	}`
	expectedSQL := `INSERT INTO swiftTable (country_iso2, code, bank_name, address, country_name, is_hq) VALUES (?, ?, ?, ?, ?, ?)`
	mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
		WithArgs("PL", "ABCDEFGH", "Polish Bank", "Polish address", "POLAND", true).
		WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest("POST", "/v1/swift-codes", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := handlers.PostCodeHandler(db)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Wrong status: received %v, expected %v", rr.Code, http.StatusCreated)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Wrong Content-Type: received %v, expected application/json", ct)
	}

	expectedResponse := `{"message":"Swift code inserted successfully"}`
	if strings.TrimSpace(rr.Body.String()) != expectedResponse {
		t.Errorf("Wrong response: received %v, expected %v", rr.Body.String(), expectedResponse)
	}

}

func TestDeleteCodeHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic("Error when creating sqlmock: " + err.Error())
	}
	defer db.Close()

	expectedSQL := `DELETE FROM swiftTable WHERE code = ?`
	mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
		WithArgs("ABCDEFGH").
		WillReturnResult(sqlmock.NewResult(0, 1))

	req := httptest.NewRequest("DELETE", "/v1/swift-codes", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("swift-code", "ABCDEFGH")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := handlers.DeleteCodeHandler(db)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Wrong status: received %v, expected %v", rr.Code, http.StatusNoContent)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Wrong Content-Type: received %v, expected application/json", ct)
	}

	expectedResponse := `{"message":"Swift code deleted successfully"}`
	if strings.TrimSpace(rr.Body.String()) != expectedResponse {
		t.Errorf("Wrong response: received %v, expected %v", rr.Body.String(), expectedResponse)
	}

}
