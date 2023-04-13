package main

import (
	"context"
	"errors"
	"flag"
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
	addr := flag.String("addr", "localhost:8080", "Server Address")
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
