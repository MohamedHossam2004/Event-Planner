package data

import (

	"go.mongodb.org/mongo-driver/bson/primitive"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventApps struct {
	EventID     primitive.ObjectID `bson:"event_id" json:"event_id"`
	Attendee []string              `bson:"attendee" json:"attendee"`
}
type EventAppModel struct {
	collection *mongo.Collection
}

func NewEventAppModel(db *mongo.Database, collectionName string) *EventAppModel {
	return &EventAppModel{
		collection: db.Collection(collectionName),
	}
}

func (s *EventAppModel) CreateEventApp(ctx context.Context, eventApp *EventApps) error {
	eventApp.EventID = primitive.NewObjectID()
	_, err := s.collection.InsertOne(ctx, eventApp)
	return err
}

func (s *EventAppModel) GetEventApp(ctx context.Context, id primitive.ObjectID) (*EventApps, error) {
	var eventApp EventApps
	err := s.collection.FindOne(ctx, bson.M{"event_id": id}).Decode(&eventApp)
	if err != nil {
		return nil, err
	}
	return &eventApp, nil
}

func (s *EventAppModel) UpdateEventApp(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := s.collection.UpdateOne(ctx, bson.M{"event_id": id}, bson.M{"$set": update})
	return err
}

func (s *EventAppModel) DeleteEventApp(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{"event_id": id})
	return err
}

func (s *EventAppModel) ListEventApps(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]*EventApps, error) {
	var eventApps []*EventApps
	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var eventApp EventApps
		if err := cursor.Decode(&eventApp); err != nil {
			return nil, err
		}
		eventApps = append(eventApps, &eventApp)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return eventApps, nil
}