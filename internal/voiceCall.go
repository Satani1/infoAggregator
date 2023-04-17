package internal

import (
	"finalWork/internal/models"
	"finalWork/internal/utils"
	"strconv"
)

func (app *Application) VoiceCall() (VoiceCallData []models.VoiceCallData, err error) {
	data, err := utils.ReadCSV("./simulator/voice.data", 8)
	if err != nil {
		return nil, err
	}
	var dataVoice [][]string
	for _, records := range data {
		if !utils.CheckCountry(records[0]) {
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
	return VoiceCallData, nil
}
