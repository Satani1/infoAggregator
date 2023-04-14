package internal

import (
	"encoding/json"
	"finalWork/internal/models"
	"net/http"
)

func (app *Application) GetResultData() (Results *models.ResultSetT, err error) {
	SMS, err := app.GetResultSMS()
	if err != nil {
		return nil, err
	}
	MMS, err := app.GetResultMMS()
	if err != nil {
		return nil, err
	}
	VoiceCall, err := app.VoiceCall()
	if err != nil {
		return nil, err
	}
	Incident, err := app.GetResultIncident()
	if err != nil {
		return nil, err
	}
	Email, err := app.GetResultEmail()
	if err != nil {
		return nil, err
	}

	Support, err := app.GetResultSupport()
	if err != nil {
		return nil, err
	}

	Billing, err := app.Billing()
	if err != nil {
		return nil, err
	}

	Results = &models.ResultSetT{
		SMS:       SMS,
		MMS:       MMS,
		VoiceCall: VoiceCall,
		Incident:  Incident,
		Email:     *Email,
		Billing:   Billing,
		Support:   Support,
	}

	return Results, nil
}

func (app *Application) GetResults(w http.ResponseWriter, r *http.Request) {
	var res models.ResultT
	Results, err := app.GetResultData()
	if (err != nil) || (Results.Support == nil || Results.Billing == nil || Results.Incident == nil || Results.MMS == nil || Results.SMS == nil || Results.VoiceCall == nil) {
		res = models.ResultT{
			Status: false,
			Data:   nil,
			Error:  "Error on collect data",
		}
	} else {
		res = models.ResultT{
			Status: true,
			Data:   Results,
			Error:  "",
		}
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "Application/json")

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)

}
