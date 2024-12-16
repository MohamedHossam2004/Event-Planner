package data

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

const dbTimeout = 3 * time.Second

type Models struct {
	Event     EventModelInterface
	EventApps EventAppModelInterface
}

func NewModels(db *mongo.Database) Models {
	eventModel := EventModel{collection: db.Collection("events")}

	return Models{
		Event: eventModel,
		EventApps: &EventAppModel{
			collection:   db.Collection("event_apps"),
			eventService: &eventModel,
		},
	}
}
