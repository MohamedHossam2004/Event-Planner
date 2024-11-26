package models

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database instance
var db *mongo.Database

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
