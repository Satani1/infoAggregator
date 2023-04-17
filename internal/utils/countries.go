package utils

import "github.com/pariz/gountries"

func CheckCountry(countryCode string) bool {
	query := gountries.New()

	_, err := query.FindCountryByAlpha(countryCode)
	if err == nil {
		return true
	}
	return false
}

func CountryName(countryCode string) (string, error) {
	query := gountries.New()

	name, err := query.FindCountryByAlpha(countryCode)
	if err != nil {
		return "", err
	}
	return name.Name.Common, nil
}
