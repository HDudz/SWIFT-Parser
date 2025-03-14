package integration_tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/HDudz/SWIFT-Parser/internal/models"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func waitForServer(url string, seconds int) error {
	var i int
	for i = 0; i < seconds; i += 2 {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {

			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		log.Print("Waiting for api to come up")
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("server not available at %s after %v", url, i)
}

func prepareDB() error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME")))

	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Truncating database tables for clean tests...")

	_, err = db.Exec("DELETE FROM swiftTable")
	if err != nil {
		return fmt.Errorf("failed to clean database: %v", err)
	}

	return nil
}

var baseURL string
var HQCode string
var BranchCode string
var ISO2Code string

func TestMain(m *testing.M) {

	baseURL = fmt.Sprintf("http://api-test:8080/v1/swift-codes")
	HQCode = "TESTCODEXXX"
	BranchCode = "TESTCODE123"
	ISO2Code = "TS"

	err := waitForServer(baseURL, 16)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to server")
	err = prepareDB()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("\n========Starting integration tests========\n\n")
	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestPostCode(t *testing.T) {

	url := baseURL
	newHQRecord := models.MainSwift{
		CountryISO2: ISO2Code,
		Code:        HQCode,
		BankName:    "Test Bank",
		Address:     "123 Test Street",
		CountryName: "TEST COUNTRY",
		IsHQ:        true,
	}

	newBranchRecord := models.MainSwift{
		CountryISO2: ISO2Code,
		Code:        BranchCode,
		BankName:    "Test Bank",
		Address:     "123 Test Street",
		CountryName: "TEST COUNTRY",
		IsHQ:        false,
	}
	payload1, err := json.Marshal(newHQRecord)
	payload2, err := json.Marshal(newBranchRecord)
	if err != nil {
		t.Fatalf("Failed to marshal new record: %v", err)
	}

	resp1, err := http.Post(url, "application/json", bytes.NewBuffer(payload1))
	resp2, err := http.Post(url, "application/json", bytes.NewBuffer(payload2))
	if err != nil {
		t.Fatalf("Failed to POST new swift code: %v", err)
	}
	defer resp1.Body.Close()
	defer resp2.Body.Close()

	if resp1.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201, got %d", resp1.StatusCode)
	} else {
		t.Logf("Correct status for resp1: %v", resp1.StatusCode)
	}
	if resp2.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201, got %d", resp2.StatusCode)
	} else {
		t.Logf("Correct status for resp2: %v", resp2.StatusCode)
	}
	body1, err := io.ReadAll(resp1.Body)
	body2, err := io.ReadAll(resp2.Body)

	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	expectedResponse := `{"message":"Swift code inserted successfully"}`

	if strings.TrimSpace(string(body1)) != expectedResponse {
		t.Errorf("Wrong response: received %v, expected %v", string(body1), expectedResponse)
	} else {
		t.Logf("Correct response for POST with code \"%s\": %s", HQCode, string(body2))
	}

	if strings.TrimSpace(string(body2)) != expectedResponse {
		t.Errorf("Wrong response: received %v, expected %v", string(body2), expectedResponse)
	} else {
		t.Logf("Correct response for POST with code \"%s\": %s", BranchCode, string(body2))
	}
}

func TestGetCode(t *testing.T) {
	url := baseURL + "/" + HQCode

	resp, err := http.Get(url)

	if err != nil {
		t.Fatalf("Failed to GET swift code: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	} else {
		t.Logf("Correct status: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	var response models.MainSwift
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("Error when decoding response: %v", err)
	}

	if response.Code != HQCode {
		t.Fatalf("Expected HQ code: %v, got \"%v\"", HQCode, response.Code)
	} else {
		t.Logf("Correct HQ code: %v", response.Code)
	}

	if len(*response.Branches) != 1 {
		t.Fatalf("Expected 1 branch connected to HQ, got \"%v\"", len(*response.Branches))
	} else {
		t.Logf("Correct number of branches connected to HQ: %v", len(*response.Branches))
	}

	if (*response.Branches)[0].Code != BranchCode {
		t.Fatalf("Expected branch code: %v, got %v", BranchCode, (*response.Branches)[0].Code)
	} else {
		t.Logf("Correct Branch code: %v", (*response.Branches)[0].Code)
	}
}

func TestGetCountry(t *testing.T) {
	url := baseURL + "/country/" + ISO2Code

	resp, err := http.Get(url)

	if err != nil {
		t.Fatalf("Failed to GET country data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	} else {
		t.Logf("Correct status: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	var response models.CountryModel
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("Error when decoding response: %v", err)
	}

	if response.CountryISO2 != ISO2Code {
		t.Fatalf("Expected ISOcode: %v, got \"%v\"", ISO2Code, response.CountryISO2)
	} else {
		t.Logf("Correct ISOcode: %v", response.CountryISO2)
	}

	if response.CountryName != "TEST COUNTRY" {
		t.Fatalf("Expected country name: TEST COUNTRY, got \"%v\"", response.CountryName)
	} else {
		t.Logf("Correct country name: %v", response.CountryName)
	}

	if len(*response.SwiftCodes) != 2 {
		t.Fatalf("Expected 2 records with \"TS\" ISO2 code, got %v", len(*response.SwiftCodes))
	} else {
		t.Logf("Correct number of records with \"TS\" ISO2 code: %v", len(*response.SwiftCodes))
	}
}

func TestDeleteCode(t *testing.T) {

	url1 := baseURL + "/" + HQCode
	url2 := baseURL + "/" + BranchCode

	client := &http.Client{}
	req1, err1 := http.NewRequest("DELETE", url1, nil)
	if err1 != nil {
		t.Fatalf("Failed to create DELETE request for url1: %v", err1)
	}
	req2, err2 := http.NewRequest("DELETE", url2, nil)
	if err2 != nil {
		t.Fatalf("Failed to create DELETE request for url2: %v", err2)
	}

	resp1, err1 := client.Do(req1)
	resp2, err2 := client.Do(req2)

	if err1 != nil {
		t.Fatalf("Failed to execute DELETE request for %v: %v", url1, err1)
	}
	if err2 != nil {
		t.Fatalf("Failed to execute DELETE request for %v: %v", url2, err2)
	}

	if resp1.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, resp1.StatusCode)
	} else {
		t.Logf("Correct status for resp1: %v", resp1.StatusCode)
	}

	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, resp2.StatusCode)
	} else {
		t.Logf("Correct status for resp2: %v", resp2.StatusCode)
	}

	body1, err := io.ReadAll(resp1.Body)
	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}
	expectedResponse := `{"message":"Swift code deleted successfully"}`

	if strings.TrimSpace(string(body1)) != expectedResponse {
		t.Errorf("Wrong response: received %v, expected %v", string(body1), expectedResponse)
	} else {
		t.Logf("Correct response for DELETE %s: %v", url1, string(body1))
	}

	if strings.TrimSpace(string(body2)) != expectedResponse {
		t.Errorf("Wrong response: received %v, expected %v", string(body2), expectedResponse)
	} else {
		t.Logf("Correct response for DELETE %s: %v", url2, string(body2))
	}

}
