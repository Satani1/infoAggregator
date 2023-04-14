package main

import (
	"context"
	"errors"
	"finalWork/internal"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	App := internal.NewApplication()

	srv := http.Server{
		Addr:     App.Addr,
		ErrorLog: App.ErrorLog,
		Handler:  App.Routes(),
	}

	App.InfoLog.Printf("Launching server on %s", App.Addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			App.ErrorLog.Fatalln(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		App.ErrorLog.Fatalln(err)
	}
	App.InfoLog.Fatalln("Im shutdown...")

}
