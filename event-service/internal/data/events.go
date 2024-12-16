package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventModelInterface interface {
	CreateEvent(event *Event) (*Event, error)
	GetEventByID(id primitive.ObjectID) (*Event, error)
	UpdateEvent(id primitive.ObjectID, event *Event) (*Event, error)
	DeleteEvent(id primitive.ObjectID) error
	GetAllEvents() ([]Event, error)
}

// EventType represents the type of event
type EventType string

const (
	Conference EventType = "CONFERENCE"
	Workshop   EventType = "WORKSHOP"
	Meetup     EventType = "MEETUP"
	Social     EventType = "SOCIAL"
	CareerFair EventType = "CAREER_FAIR"
	Graduation EventType = "GRADUATION"
	Other      EventType = "OTHER"
)

type EventModel struct {
	collection *mongo.Collection
}

// Event represents the main event document structure
type Event struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Date                 time.Time          `bson:"date" json:"date" validate:"required"`
	Type                 EventType          `bson:"type" json:"type" validate:"required"`
	Name                 string             `bson:"name" json:"name" validate:"required"`
	Location             Location           `bson:"location" json:"location" validate:"required"`
	NumberOfApplications int                `bson:"number_of_applications" json:"number_of_applications"`
	Ushers               []string           `bson:"ushers" json:"ushers"`
	Description          string             `bson:"description" json:"description" validate:"required"`
	MaxCapacity          int                `bson:"max_capacity" json:"max_capacity" validate:"required"`
	MinCapacity          int                `bson:"min_capacity" json:"min_capacity" validate:"required"`
	Organizers           []Organizer        `bson:"organizers" json:"organizers" validate:"required,min=1"`
	CreatedAt            time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt            time.Time          `bson:"updated_at" json:"updated_at"`
	Status               string             `bson:"status" json:"status"`
}

// Location represents the event location details
type Location struct {
	Address string `bson:"address" json:"address" validate:"required"`
	City    string `bson:"city" json:"city" validate:"required"`
	State   string `bson:"state" json:"state" validate:"required"`
	Country string `bson:"country" json:"country" validate:"required"`
}

// Organizer represents the event organizer details
type Organizer struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name  string             `bson:"name" json:"name" validate:"required"`
	Email string             `bson:"email" json:"email" validate:"required,email"`
	Phone string             `bson:"phone" json:"phone"`
	Role  string             `bson:"role" json:"role"`
}

// CreateIndexes creates the necessary indexes for the Event collection
func CreateEventIndexes() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "date", Value: 1},
				{Key: "type", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "name", Value: "text"},
			},
		},
		{
			Keys: bson.D{
				{Key: "location.city", Value: 1},
				{Key: "location.country", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
			},
		},
	}
}

// NewEventService creates a new instance of EventService
func NewEventService(db *mongo.Database, collectionName string) *EventModel {
	return &EventModel{
		collection: db.Collection(collectionName),
	}
}

// CreateEvent adds a new event to the database
func (es EventModel) CreateEvent(event *Event) (*Event, error) {
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
func (es EventModel) GetEventByID(id primitive.ObjectID) (*Event, error) {
	var event Event
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
func (es EventModel) UpdateEvent(id primitive.ObjectID, event *Event) (*Event, error) {
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
func (es EventModel) DeleteEvent(id primitive.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: id}}

	_, err := es.collection.DeleteOne(context.Background(), filter)
	return err
}

// GetAllEvents retrieves all events
func (es EventModel) GetAllEvents() ([]Event, error) {
	var events []Event
	cursor, err := es.collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var event Event
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
