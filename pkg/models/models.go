package models

type SMSData struct {
	Country      string
	Bandwidth    string
	ResponseTime string
	Provider     string
}

type MMSData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	ResponseTime string `json:"response_time"`
	Bandwidth    string `json:"bandwidth"`
}
