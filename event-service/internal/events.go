package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// EventType represents the type of event
type EventType string

const (
	Conference EventType = "CONFERENCE"
	Workshop   EventType = "WORKSHOP"
	Meetup     EventType = "MEETUP"
	Social     EventType = "SOCIAL"
	Other      EventType = "OTHER"
)

// Event represents the main event document structure
type Event struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Date        time.Time          `bson:"date" json:"date" validate:"required"`
	Type        EventType          `bson:"type" json:"type" validate:"required"`
	Name        string             `bson:"name" json:"name" validate:"required"`
	Location    Location           `bson:"location" json:"location" validate:"required"`
	NumberOfApplications     int                `bson:"number_of_applications" json:"number_of_applications"`
	MaxCapacity int                `bson:"max_capacity" json:"max_capacity" validate:"required"`
	MinCapacity int                `bson:"min_capacity" json:"min_capacity" validate:"required"`
	Organizers  []Organizer        `bson:"organizers" json:"organizers" validate:"required,min=1"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	Status      string             `bson:"status" json:"status"`
}

// Location represents the event location details
type Location struct {
	Address     string    `bson:"address" json:"address" validate:"required"`
	City        string    `bson:"city" json:"city" validate:"required"`
	State       string    `bson:"state" json:"state" validate:"required"`
	Country     string    `bson:"country" json:"country" validate:"required"`
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
