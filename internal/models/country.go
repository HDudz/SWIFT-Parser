package models

type CountryModel struct {
	CountryISO2 string      `json:"countryISO2"`
	CountryName string      `json:"countryName"`
	SwiftCodes  *[]SubSwift `json:"swiftCodes,omitempty"`
}
