package internal

import (
	"encoding/json"
	"errors"
	"finalWork/internal/models"
	"finalWork/internal/utils"
)

func (app *Application) Support() (SupportData []models.SupportData, err error) {
	requestURL := "http://127.0.0.1:8383/support"
	resBody, err := utils.SendGetRequest(requestURL)
	if err != nil {
		return nil, err
	} else if errors.Is(err, errors.New("status code is 500")) {
		return nil, err
	}

	if err := json.Unmarshal(resBody, &SupportData); err != nil {
		return SupportData, err
	}
	return SupportData, nil
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
