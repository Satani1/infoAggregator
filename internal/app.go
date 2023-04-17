package internal

import (
	"flag"
	"log"
	"os"
)

type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Addr     string
}

func NewApplication() *Application {
	//addr config from terminal
	addr := flag.String("addr", "localhost:8080", "Server Address")
	flag.Parse()
	//logs
	InfoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	ErrorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	App := &Application{
		ErrorLog: ErrorLog,
		InfoLog:  InfoLog,
		Addr:     *addr,
	}
	return App
}
