package models

type SwiftModel struct {
	ID          int    `json:"id"`
	CountryISO2 string `json:"countryISO2"`
	Code        string `json:"swiftCode"`
	CodeType    string `json:"codeType"`
	BankName    string `json:"bankName"`
	Address     string `json:"address"`
	Town        string `json:"town"`
	CountryName string `json:"countryName"`
	TimeZone    string `json:"timeZone"`
	IsHQ        bool   `json:"isHeadquarter"`
}
