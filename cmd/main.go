package main

import (
	"context"
	"errors"
	"finalWork/pkg"
	"finalWork/pkg/models"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	Addr     string
}

func NewApplication() *Application {
	//addr config from terminal
	addr := flag.String("addr", "localhost:8282", "Server Address")
	flag.Parse()
	//logs
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	App := &Application{
		errorLog: errorLog,
		infoLog:  infoLog,
		Addr:     *addr,
	}
	return App
}

func main() {

	App := NewApplication()

	srv := http.Server{
		Addr:     App.Addr,
		ErrorLog: App.errorLog,
		Handler:  App.Routes(),
	}
	SMS()
	App.infoLog.Printf("Launching server on %s", App.Addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			App.errorLog.Fatalln(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		App.errorLog.Fatalln(err)
	}
	App.infoLog.Fatalln("Im shutdown...")

}

var Alpha2Countries = map[string]string{
	"RU": "Russia",
	"AF": "Afghanistan",
}

func SMS() (SMSData []models.SMSData, err error) {
	data, err := pkg.ReadCSV("./data/sms.data", 4)
	if err != nil {
		return nil, err
	}
	fmt.Println(data)
	var dataSMS [][]string
	for _, records := range data {
		if _, found := Alpha2Countries[records[0]]; !found {
			break
		}
		if !(records[3] == "Topolo" || records[3] == "Rond" || records[3] == "Kildy") {
			break
		}
		dataSMS = append(dataSMS, records)
	}

	for _, records := range data {
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
