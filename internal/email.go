package internal

import (
	"finalWork/internal/models"
	"finalWork/internal/utils"
	"sort"
	"strconv"
)

func (app *Application) Email() (EmailData []models.EmailData, err error) {
	data, err := utils.ReadCSV("./simulator/email.data", 3)
	if err != nil {
		return nil, err
	}
	var dataEmail [][]string
	providers := []string{"Gmail", "Yahoo", "Hotmail", "MSN", "Orange", "Comcast", "AOL", "Live", "RediffMail", "GMX", "Protonmail",
		"Yandex", "Mail.ru"}
	for _, records := range data {
		var found bool
		for _, provider := range providers {
			found = false
			if records[1] == provider {
				found = true
			}
		}
		if found {
			continue
		}
		if !utils.CheckCountry(records[0]) {
			continue
		}
		dataEmail = append(dataEmail, records)
	}
	for _, records := range dataEmail {
		dTime, err := strconv.Atoi(records[2])
		if err != nil {
			return nil, err
		}
		var email = &models.EmailData{
			Country:      records[0],
			Provider:     records[1],
			DeliveryTime: dTime,
		}

		EmailData = append(EmailData, *email)
	}
	return EmailData, nil
}

func (app *Application) GetResultEmail() (*map[string][][]models.EmailData, error) {
	EmailData, err := app.Email()
	if err != nil {
		return nil, err
	}
	EmailDataByCountries := EmailData

	sort.SliceStable(EmailDataByCountries, func(i, j int) bool {
		return EmailDataByCountries[i].Country < EmailDataByCountries[j].Country
	})
	EmailResult := make(map[string][][]models.EmailData, 0)

	tempArr := make([]models.EmailData, 0)

	for i := 0; i < len(EmailDataByCountries)-1; i++ {
		tempCountry, err := utils.CountryName(EmailDataByCountries[i].Country)
		if err != nil {
			return nil, err
		}
		tempCountry2, err := utils.CountryName(EmailDataByCountries[i+1].Country)
		if err != nil {
			return nil, err
		}
		if tempCountry == tempCountry2 {
			tempArr = append(tempArr, EmailDataByCountries[i])
		} else {
			sort.SliceStable(tempArr, func(i, j int) bool {
				return tempArr[i].DeliveryTime < tempArr[j].DeliveryTime
			})
			tempArrH := tempArr[0:3]
			tempArrL := tempArr[len(tempArr)-3 : len(tempArr)]
			arr := [][]models.EmailData{tempArrH, tempArrL}
			EmailResult[tempCountry] = arr
			tempArrH, tempArrL, tempArr = nil, nil, nil
		}
	}

	return &EmailResult, nil
}
