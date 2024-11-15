package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const webPort = "80"

type application struct{}

func main() {
	app := &application{}
	log.Printf("starting broker service on %s\n", webPort)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", webPort),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := srv.ListenAndServe()
	log.Fatal(err)
}
