package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/MohamedHossam2004/Event-Planner/user-service/internal/data"
	"github.com/MohamedHossam2004/Event-Planner/user-service/internal/mailer"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var counts int64

const (
	webPort = "80"
	webEnv  = "development"
)

type config struct {
	port string
	env  string
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	jwt struct {
		secret string
	}
}

type application struct {
	config config
	logger *log.Logger
	models data.Models
	mailer mailer.Mailer
}

func main() {
	var cfg config
	cfg.port = webPort
	cfg.env = webEnv
	cfg.smtp.host = os.Getenv("MAILHOG_HOST")
	portStr := os.Getenv("MAILHOG_PORT")
	cfg.smtp.username = os.Getenv("MAILHOG_USERNAME")
	cfg.smtp.password = os.Getenv("MAILHOG_PASSWORD")
	cfg.smtp.sender = os.Getenv("SENDER_EMAIL")
	cfg.jwt.secret = os.Getenv("JWT_SECRET")

	if cfg.smtp.host == "" || portStr == "" {
		log.Fatal("Environment variables for Mailhog are not set")
	}
	var err error
	cfg.smtp.port, err = strconv.Atoi(portStr)
	if err != nil {
		log.Fatal("Error: Invalid MAILHOG_PORT value: %s\n", portStr)
		return
	}

	db := connectToDB()
	if db == nil {
		log.Panic("could not connect to database")
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.LUTC)
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
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
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Database not yet ready...")
			counts++
		} else {
			log.Println("Connected to Database!")
			return connection
		}
		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing of for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
