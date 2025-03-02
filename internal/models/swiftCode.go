package models

type SubModel struct {
	Address     string `json:"address"`
	BankName    string `json:"bankName"`
	CountryISO2 string `json:"countryISO2"`
	IsHQ        bool   `json:"isHeadquarter"`
	Code        string `json:"swiftCode"`
}
type MainModel struct {
	Address     string      `json:"address"`
	BankName    string      `json:"bankName"`
	CountryISO2 string      `json:"countryISO2"`
	CountryName string      `json:"countryName"`
	IsHQ        bool        `json:"isHeadquarter"`
	Code        string      `json:"swiftCode"`
	Branches    *[]SubModel `json:"branches,omitempty"`
}
