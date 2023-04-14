package internal

import (
	"encoding/json"
	"errors"
	"finalWork/internal/models"
	"finalWork/internal/utils"
	"sort"
)

func (app *Application) MMS() (MMSData []models.MMSData, err error) {
	requestURL := "http://127.0.0.1:8383/mms"
	resBody, err := utils.SendGetRequest(requestURL)
	if err != nil {
		return nil, err
	} else if errors.Is(err, errors.New("status code is 500")) {
		return MMSData, err
	}

	var MMSDataSlice []models.MMSData

	if err := json.Unmarshal(resBody, &MMSDataSlice); err != nil {
		return MMSData, err
	}

	for _, element := range MMSDataSlice {
		if !utils.CheckCountry(element.Country) {
			continue
		}
		if !(element.Provider == "Topolo" || element.Provider == "Rond" || element.Provider == "Kildy") {
			continue
		}
		MMSData = append(MMSData, element)

	}
	return MMSData, nil
}

func (app *Application) GetResultMMS() (MMSDataResult [][]models.MMSData, err error) {
	MMSData, err := app.MMS()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(MMSData); i++ {
		countryName, err := utils.CountryName(MMSData[i].Country)
		if err != nil {
			return nil, err
		}
		MMSData[i].Country = countryName

	}

	SMSDataByCountry := MMSData
	SMSDataByProvider := MMSData

	sort.SliceStable(SMSDataByCountry, func(i, j int) bool {
		return SMSDataByCountry[i].Country < SMSDataByCountry[j].Country
	})

	sort.SliceStable(SMSDataByProvider, func(i, j int) bool {
		return SMSDataByProvider[i].Provider < SMSDataByProvider[j].Provider
	})

	MMSDataResult = [][]models.MMSData{SMSDataByCountry, SMSDataByProvider}

	return MMSDataResult, nil
}
