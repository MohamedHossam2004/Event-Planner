package data

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventApps struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	EventID   primitive.ObjectID `bson:"event_id" json:"event_id" validate:"required"`
	Attendee  []string          `bson:"attendee" json:"attendee" validate:"required"`
}

type EventAppModel struct {
	collection   *mongo.Collection
	eventService *EventModel
}

func NewEventAppModel(db *mongo.Database, collectionName string, eventService *EventModel) *EventAppModel {
	return &EventAppModel{
		collection:   db.Collection(collectionName),
		eventService: eventService,
	}
}

func (s *EventAppModel) CreateEventApp(ctx context.Context, eventApp *EventApps) error {
	// Verify that the event exists
	event, err := s.eventService.GetEventByID(eventApp.EventID)
	if err != nil {
		return err
	}
	if event == nil {
		return fmt.Errorf("event with ID %s does not exist", eventApp.EventID.Hex())
	}

	eventApp.ID = primitive.NewObjectID()
	
	// Insert the event application
	_, err = s.collection.InsertOne(ctx, eventApp)
	if err != nil {
		return err
	}

	// Update the number_of_applications in the Event document
	update := bson.M{"$inc": bson.M{"number_of_applications": 1}}
	_, err = s.eventService.collection.UpdateOne(
		ctx,
		bson.M{"_id": eventApp.EventID},
		update,
	)
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