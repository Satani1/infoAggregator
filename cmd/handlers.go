package main

import (
	"encoding/json"
	"errors"
	"finalWork/pkg"
	"finalWork/pkg/models"
	"fmt"
	"github.com/pariz/gountries"
	"io"
	"net/http"
	"strconv"
)

func CheckCountry(countryCode string) bool {
	query := gountries.New()

	_, err := query.FindCountryByAlpha(countryCode)
	if err == nil {
		return true
	}
	return false
}

func (app *Application) SendGetRequest(requestURL string) (resBody []byte, err error) {
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New("status code is " + strconv.Itoa(res.StatusCode))
	}

	resBody, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return resBody, nil

}

func SMS() (SMSData []models.SMSData, err error) {
	data, err := pkg.ReadCSV("./data/sms.data", 4)
	if err != nil {
		return nil, err
	}
	fmt.Println(data)
	var dataSMS [][]string
	for _, records := range data {
		if !CheckCountry(records[0]) {
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
	fmt.Println(SMSData)
	return SMSData, nil
}

func (app *Application) MMS() (MMSData []models.MMSData, err error) {
	requestURL := "http://127.0.0.1:8383/mms"
	resBody, err := app.SendGetRequest(requestURL)
	if err != nil {
		return nil, err
	} else if errors.Is(err, errors.New("status code is 500")) {
		return MMSData, err
	}
	fmt.Printf("%s\n", resBody)

	var MMSDataSlice []models.MMSData

	if err := json.Unmarshal(resBody, &MMSDataSlice); err != nil {
		return MMSData, err
	}
	fmt.Println(MMSDataSlice)

	for _, element := range MMSDataSlice {
		if !CheckCountry(element.Country) {
			continue
		}
		if !(element.Provider == "Topolo" || element.Provider == "Rond" || element.Provider == "Kildy") {
			continue
		}
		MMSData = append(MMSData, element)

	}
	fmt.Println(MMSData)
	return MMSData, nil
}
