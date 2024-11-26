package data

import (
	"time"
	"go.mongodb.org/mongo-driver/mongo"
)
const dbTimeout = 3 * time.Second
type Models struct {
	Event EventModel
	EventApps EventAppModel
}
func NewModels(db *mongo.Database) Models {
	return Models{
		Event: EventModel{collection: db.Collection("events")},
		EventApps: EventAppModel{collection: db.Collection("event_apps")},
	}
}