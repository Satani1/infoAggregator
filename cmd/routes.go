package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func (app *Application) Routes() *mux.Router {
	rMux := mux.NewRouter()

	rMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello, im service for final work. My addres is - " + app.Addr)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	rMux.HandleFunc("/mms", func(w http.ResponseWriter, r *http.Request) {
		out, err := app.MMS()
		if err != nil {
			app.errorLog.Fatalln(err)
		}
		fmt.Println(out)
	})

	rMux.HandleFunc("/support", func(w http.ResponseWriter, r *http.Request) {
		out, err := app.Support()
		if err != nil {
			app.errorLog.Fatalln(err)
		}
		fmt.Println(out)
	})

	rMux.HandleFunc("/accendent", func(w http.ResponseWriter, r *http.Request) {
		out, err := app.Incident()
		if err != nil {
			app.errorLog.Fatalln(err)
		}
		fmt.Println(out)
	})

	rMux.HandleFunc("/api", app.GetResults)

	return rMux
}
