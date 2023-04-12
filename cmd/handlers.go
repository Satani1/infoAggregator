package main

import (
	"encoding/json"
	"errors"
	"finalWork/pkg"
	"finalWork/pkg/models"
	"fmt"
	"github.com/pariz/gountries"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
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

func CountryName(countryCode string) (string, error) {
	query := gountries.New()

	name, err := query.FindCountryByAlpha(countryCode)
	if err != nil {
		return "", err
	}
	return name.Name.Common, nil
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

func (app *Application) SMS() (SMSData []models.SMSData, err error) {
	data, err := pkg.ReadCSV("./data/sms.data", 4)
	if err != nil {
		return nil, err
	}
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

func (app *Application) VoiceCall() (VoiceCallData []models.VoiceCallData, err error) {
	data, err := pkg.ReadCSV("./data/voice.data", 8)
	if err != nil {
		return nil, err
	}
	fmt.Println(data)
	var dataVoice [][]string
	for _, records := range data {
		if !CheckCountry(records[0]) {
			continue
		}
		if !(records[3] == "TransparentCalls" || records[3] == "E-Voice" || records[3] == "JustPhone") {
			continue
		}
		dataVoice = append(dataVoice, records)
	}

	for _, records := range dataVoice {
		ConnStab, err := strconv.ParseFloat(records[4], 32)
		if err != nil {
			return nil, err
		}
		TimeToFirstByte, err := strconv.Atoi(records[5])
		if err != nil {
			return nil, err
		}
		VoicePur, err := strconv.Atoi(records[6])
		if err != nil {
			return nil, err
		}
		MedianTime, err := strconv.Atoi(records[7])
		if err != nil {
			return nil, err
		}
		var voice = &models.VoiceCallData{
			Country:             records[0],
			Bandwidth:           records[1],
			ResponseTime:        records[2],
			Provider:            records[3],
			ConnectionStability: float32(ConnStab),
			TTFB:                TimeToFirstByte,
			VoicePurity:         VoicePur,
			MedianOfCallsTime:   MedianTime,
		}

		VoiceCallData = append(VoiceCallData, *voice)
	}
	fmt.Println(VoiceCallData)
	return VoiceCallData, nil
}

func (app *Application) Email() (EmailData []models.EmailData, err error) {
	data, err := pkg.ReadCSV("./data/email.data", 3)
	if err != nil {
		return nil, err
	}
	fmt.Println(data)
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
		if !CheckCountry(records[0]) {
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
	fmt.Println(EmailData)
	return EmailData, nil
}

func (app *Application) Billing() (BillingData models.BillingData, err error) {
	nums, err := ioutil.ReadFile("./data/billing.data")
	if err != nil {
		return BillingData, err
	}
	//fmt.Printf("%v, %v", nums, reflect.TypeOf(nums))
	var Bits []int
	for i := 5; i >= 0; i-- {
		byteNumber, _ := strconv.Atoi(string(nums[i]))
		Bits = append(Bits, byteNumber)
	}
	var sum uint8
	for index, num := range Bits {
		if num == 1 {
			sum += uint8(math.Pow(float64(2), float64(index)))
		}
	}

	BillingData = models.BillingData{

		CreateCustomer: BillCheck(nums[0]),
		Purchase:       BillCheck(nums[1]),
		Payout:         BillCheck(nums[2]),
		Recurring:      BillCheck(nums[3]),
		FraudControl:   BillCheck(nums[4]),
		CheckoutPage:   BillCheck(nums[5]),
	}

	return BillingData, nil
}
func BillCheck(bit uint8) bool {
	if int(bit) == 49 {
		return true
	} else {
		return false
	}
}

func (app *Application) Support() (SupportData []models.SupportData, err error) {
	requestURL := "http://127.0.0.1:8383/support"
	resBody, err := app.SendGetRequest(requestURL)
	if err != nil {
		return nil, err
	} else if errors.Is(err, errors.New("status code is 500")) {
		return nil, err
	}

	if err := json.Unmarshal(resBody, &SupportData); err != nil {
		return SupportData, err
	}
	fmt.Println(SupportData)
	return SupportData, nil
}

func (app *Application) Incident() (IncidentData []models.IncidentData, err error) {
	requestURL := "http://localhost:8383/accendent"
	resBody, err := app.SendGetRequest(requestURL)
	if err != nil {
		return nil, err
	} else if errors.Is(err, errors.New("status code is 500")) {
		return nil, err
	}

	if err := json.Unmarshal(resBody, &IncidentData); err != nil {
		return nil, err
	}
	fmt.Println(IncidentData)
	return IncidentData, nil
}

func (app *Application) GetResultSMS() (SMSDataResult [][]models.SMSData, err error) {
	SMSData, err := app.SMS()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(SMSData); i++ {
		countryName, err := CountryName(SMSData[i].Country)
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

func (app *Application) GetResultMMS() (MMSDataResult [][]models.MMSData, err error) {
	MMSData, err := app.MMS()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(MMSData); i++ {
		countryName, err := CountryName(MMSData[i].Country)
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

func (app *Application) GetResultEmail() (map[string][][]models.EmailData, error) {
	EmailData, err := app.Email()
	if err != nil {
		return nil, err
	}
	EmailDataHigh := EmailData
	EmailDataLow := EmailData
	sort.SliceStable(EmailDataHigh, func(i, j int) bool {
		return EmailDataHigh[i].DeliveryTime < EmailDataHigh[j].DeliveryTime
	})
	sort.SliceStable(EmailDataLow, func(i, j int) bool {
		return EmailDataLow[i].DeliveryTime > EmailDataLow[j].DeliveryTime
	})
	EmailResult := make(map[string][][]models.EmailData, 0)

	return EmailResult, nil
}

func (app *Application) GetResultIncident() (IncidentResult []models.IncidentData, err error) {
	IncidentData, err := app.Incident()
	if err != nil {
		return nil, err
	}
	sort.SliceStable(IncidentData, func(i, j int) bool {
		return IncidentData[i].Status < IncidentData[j].Status
	})

	return IncidentData, nil
}

func (app *Application) GetResultSupport() (data []int, err error) {
	SupportData, err := app.Support()
	if err != nil {
		return nil, err
	}

	var sum int
	for _, ticket := range SupportData {
		sum += ticket.ActiveTickets
	}

	if sum < 9 {
		data = append(data, 1)
	} else if sum < 16 {
		data = append(data, 2)
	} else {
		data = append(data, 3)
	}

	averageTime := 60 / 18
	data = append(data, sum*averageTime)

	return data, nil
}

func (app *Application) GetResultData() (Results models.ResultSetT) {
	SMS, err := app.GetResultSMS()
	if err != nil {
		app.errorLog.Fatalln(err)
	}
	MMS, err := app.GetResultMMS()
	if err != nil {
		app.errorLog.Fatalln(err)
	}
	VoiceCall, err := app.VoiceCall()
	if err != nil {
		app.errorLog.Fatalln(err)
	}
	Incident, err := app.GetResultIncident()
	if err != nil {
		app.errorLog.Fatalln(err)
	}
	//Email, err := app.GetResultEmail()
	//if err != nil {
	//	app.errorLog.Fatalln(err)
	//}

	Support, err := app.GetResultSupport()
	if err != nil {
		app.errorLog.Fatalln(err)
	}

	Billing, err := app.Billing()
	if err != nil {
		app.errorLog.Fatalln(err)
	}
	if err == nil {
		Results = models.ResultSetT{
			SMS:       SMS,
			MMS:       MMS,
			VoiceCall: VoiceCall,
			Incident:  Incident,
			//Email:     Email,
			Billing: Billing,
			Support: Support,
		}
	} else {
		return
	}
	return Results
}

func (app *Application) GetResults(w http.ResponseWriter, r *http.Request) {
	Results := app.GetResultData()
	var res models.ResultT
	res = models.ResultT{
		Status: true,
		Data:   Results,
		Error:  "",
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
