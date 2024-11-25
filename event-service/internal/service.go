// event-service.go

package models

import (
	"context"
	"time"
	 // Adjust the import path to your project structure
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

// EventService holds the methods related to event operations
type EventService struct {
	collection *mongo.Collection
}

// NewEventService creates a new instance of EventService
func NewEventService(db *mongo.Database, collectionName string) *EventService {
	return &EventService{
		collection: db.Collection(collectionName),
	}
}

// CreateEvent adds a new event to the database
func (es *EventService) CreateEvent(event *models.Event) (*models.Event, error) {
	event.ID = primitive.NewObjectID()
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()
	event.Status = "PENDING" // Initial status of an event

	_, err := es.collection.InsertOne(context.Background(), event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// GetEventByID retrieves an event by its ID
func (es *EventService) GetEventByID(id primitive.ObjectID) (*models.Event, error) {
	var event models.Event
	filter := bson.D{{Key: "_id", Value: id}}

	err := es.collection.FindOne(context.Background(), filter).Decode(&event)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No event found
		}
		return nil, err
	}

	return &event, nil
}

// UpdateEvent updates an existing event
func (es *EventService) UpdateEvent(id primitive.ObjectID, event *models.Event) (*models.Event, error) {
	event.UpdatedAt = time.Now()
	filter := bson.D{{Key: "_id", Value: id}}

	update := bson.D{
		{Key: "$set", Value: event},
	}

	_, err := es.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	return es.GetEventByID(id)
}

// DeleteEvent removes an event from the database
func (es *EventService) DeleteEvent(id primitive.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: id}}

	_, err := es.collection.DeleteOne(context.Background(), filter)
	return err
}

// GetAllEvents retrieves all events
func (es *EventService) GetAllEvents() ([]models.Event, error) {
	var events []models.Event
	cursor, err := es.collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var event models.Event
		if err := cursor.Decode(&event); err != nil {
			log.Println("Error decoding event:", err)
			continue
		}
		events = append(events, event)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
