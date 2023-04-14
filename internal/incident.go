package internal

import (
	"encoding/json"
	"errors"
	"finalWork/internal/models"
	"finalWork/internal/utils"
	"sort"
)

func (app *Application) Incident() (IncidentData []models.IncidentData, err error) {
	requestURL := "http://localhost:8383/accendent"
	resBody, err := utils.SendGetRequest(requestURL)
	if err != nil {
		return nil, err
	} else if errors.Is(err, errors.New("status code is 500")) {
		return nil, err
	}

	if err := json.Unmarshal(resBody, &IncidentData); err != nil {
		return nil, err
	}
	return IncidentData, nil
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
