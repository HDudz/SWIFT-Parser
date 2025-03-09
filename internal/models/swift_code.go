package models

type SubSwift struct {
	Address     string `json:"address"`
	BankName    string `json:"bankName"`
	CountryISO2 string `json:"countryISO2"`
	IsHQ        bool   `json:"isHeadquarter"`
	Code        string `json:"swiftCode"`
}
type MainSwift struct {
	Address     string      `json:"address"`
	BankName    string      `json:"bankName"`
	CountryISO2 string      `json:"countryISO2"`
	CountryName string      `json:"countryName"`
	IsHQ        bool        `json:"isHeadquarter"`
	Code        string      `json:"swiftCode"`
	Branches    *[]SubSwift `json:"branches,omitempty"`
}
