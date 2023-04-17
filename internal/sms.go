package internal

import (
	"finalWork/internal/models"
	"finalWork/internal/utils"
	"sort"
)

func (app *Application) SMS() (SMSData []models.SMSData, err error) {
	data, err := utils.ReadCSV("./simulator/sms.data", 4)
	if err != nil {
		return nil, err
	}
	var dataSMS [][]string
	for _, records := range data {
		if !utils.CheckCountry(records[0]) {
			continue
		}
		if !(records[3] == "Topolo" || records[3] == "Rond" || records[3] == "Kildy") {
			continue
		}
		dataSMS = append(dataSMS, records)
	}

	for _, records := range dataSMS {
		var sms = &models.SMSData{
			Country:      records[0],
			Bandwidth:    records[1],
			ResponseTime: records[2],
			Provider:     records[3],
		}

		SMSData = append(SMSData, *sms)
	}
	return SMSData, nil
}

func (app *Application) GetResultSMS() (SMSDataResult [][]models.SMSData, err error) {
	SMSData, err := app.SMS()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(SMSData); i++ {
		countryName, err := utils.CountryName(SMSData[i].Country)
		if err != nil {
			return nil, err
		}
		SMSData[i].Country = countryName

	}

	SMSDataByCountry := SMSData
	SMSDataByProvider := SMSData

	sort.SliceStable(SMSDataByCountry, func(i, j int) bool {
		return SMSDataByCountry[i].Country < SMSDataByCountry[j].Country
	})

	sort.SliceStable(SMSDataByProvider, func(i, j int) bool {
		return SMSDataByProvider[i].Provider < SMSDataByProvider[j].Provider
	})

	SMSDataResult = [][]models.SMSData{SMSDataByCountry, SMSDataByProvider}
	return SMSDataResult, nil

}
