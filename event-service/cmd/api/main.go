package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/MohamedHossam2004/Event-Planner/event-service/internal/data"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database instance
var db *mongo.Database

const (
	webPort  = "80"
	webEnv   = "development"
	mongoURL = "mongodb://mongo:27017"
)

type config struct {
	port string
	env  string
	jwt  struct {
		secret string
	}
}

type application struct {
	config         config
	Logger         *log.Logger
	models         data.Models
	Rabbit         *amqp.Connection
	tokenExtractor TokenExtractor
}

func main() {
	var cfg config
	cfg.port = webPort
	cfg.env = webEnv
	cfg.jwt.secret = os.Getenv("JWT_SECRET")

	// Connect to the MongoDB database
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	db = mongoClient.Database("events")
	if db == nil {
		log.Panic("could not connect to database")
	}

	// Connect to RabbitMQ
	rabbitConn, err := connectToRabbit()
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitConn.Close()
	log.Println("Connected to RabbitMQ")

	// Create a logger
	logger := log.New(os.Stdout, "", log.Ldate|log.LUTC)

	// Initialize the application with models and configuration
	app := &application{
		config: cfg,
		Logger: logger,
		models: data.NewModels(db),
		Rabbit: rabbitConn,
		tokenExtractor: &realTokenExtractor{
			jwtSecret: cfg.jwt.secret,
		},
	}

	// Log the server start
	log.Printf("starting events service on %s\n", cfg.port)

	// Set up HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.port),
		Handler:      app.routes(), // Routing goes here
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start the server and log any errors
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	log.Fatal(err)
}

// Connect initializes the database connection for the Event-service
func connect() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Connect to MongoDB (adjust URI if needed)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Access the Event-service database
	db = client.Database("Event-service")
	log.Println("Connected to MongoDB, using database Event-service")
}

func connectToMongo() (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	log.Println("Connected to mongo!")

	return c, nil
}

// GetCollection returns a specific collection
func GetCollection(collectionName string) *mongo.Collection {
	return db.Collection(collectionName)
}

func connectToRabbit() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ Not yet ready...")
			counts++
		} else {
			connection = c
			break
		}
		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Printf("Retrying in %v\n", backOff)
		time.Sleep(backOff)
	}

	return connection, nil
}
