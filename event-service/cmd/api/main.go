package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MohamedHossam2004/Event-Planner/event-service/internal/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database instance
var db *mongo.Database
const (
	webPort = "80"
	webEnv  = "development"
)

type config struct {
	port string
	env  string
}

type application struct {
	config config
	Logger *log.Logger
	models data.Models
}

func main() {
	var cfg config
	cfg.port = webPort
	cfg.env = webEnv

	// Connect to the MongoDB database
	Connect()
	if db == nil {
		log.Panic("could not connect to database")
	}
	

	// Create a logger
	logger := log.New(os.Stdout, "", log.Ldate|log.LUTC)

	// Initialize the application with models and configuration
	app := &application{
		config: cfg,
		Logger: logger,
		models: data.NewModels(db),
	}

	// Log the server start
	log.Printf("starting user service on %s\n", cfg.port)

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
	err := srv.ListenAndServe()
	log.Fatal(err)
}

// Connect initializes the database connection for the Event-service
func Connect() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Connect to MongoDB (adjust URI if needed)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Access the Event-service database
	db = client.Database("Event-service")
	log.Println("Connected to MongoDB, using database Event-service")
}

// GetCollection returns a specific collection
func GetCollection(collectionName string) *mongo.Collection {
	return db.Collection(collectionName)
}
