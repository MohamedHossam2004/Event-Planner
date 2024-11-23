package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var counts int64

const (
	webPort = "80"
	webEnv  = "development"
)

type config struct {
	port string
	env  string
	jwt  struct {
		secret string
	}
}

type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config
	cfg.port = webPort
	cfg.env = webEnv
	cfg.jwt.secret = os.Getenv("JWT_SECRET")

	logger := log.New(os.Stdout, "", log.Ldate|log.LUTC)
	app := &application{
		config: cfg,
		logger: logger,
	}

	log.Printf("starting user service on %s\n", cfg.port)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	log.Fatal(err)
}
