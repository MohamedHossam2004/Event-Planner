package data

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrEventEnded     = errors.New("Event is finished")
	ErrAlreadyApplied = errors.New("User has already applied for this event")
	ErrNotApplied     = errors.New("User didn't apply to event")
)

type EventAppModelInterface interface {
	CreateEventApp(ctx context.Context, eventApp *EventApps) error
	GetEventApp(ctx context.Context, id primitive.ObjectID) (*EventApps, error)
	UpdateEventApp(ctx context.Context, id primitive.ObjectID, update bson.M) error
	DeleteEventApp(ctx context.Context, id primitive.ObjectID) error
	ListEventApps(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]*EventApps, error)
	AddAttendeeToEvent(name string, eventId primitive.ObjectID) error
	RemoveAttendeeFromEvent(name string, eventId primitive.ObjectID) error
	GetEventsByUserEmail(email string) ([]*Event, error)
}

type EventApps struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	EventID  primitive.ObjectID `bson:"event_id" json:"event_id" validate:"required"`
	Attendee []string           `bson:"attendee" json:"attendee" validate:"required"`
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

	return err
}

func (s *EventAppModel) GetEventApp(ctx context.Context, id primitive.ObjectID) (*EventApps, error) {
	var eventApp EventApps
	err := s.collection.FindOne(ctx, bson.M{"event_id": id}).Decode(&eventApp)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoRecords
		}
		return &EventApps{}, err
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

func (e *EventAppModel) AddAttendeeToEvent(name string, eventId primitive.ObjectID) error {
	_, err := e.collection.UpdateOne(context.Background(), bson.M{"event_id": eventId}, bson.M{"$push": bson.M{"attendee": name}})
	if err != nil {
		return err
	}
	update := bson.M{"$inc": bson.M{"number_of_applications": -1}}
	_, err = e.eventService.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": eventId},
		update,
	)
	return err
}

func (e *EventAppModel) RemoveAttendeeFromEvent(name string, eventId primitive.ObjectID) error {
	_, err := e.collection.UpdateOne(context.Background(), bson.M{"event_id": eventId}, bson.M{"$pull": bson.M{"attendee": name}})
	if err != nil {
		return err
	}
	update := bson.M{"$inc": bson.M{"number_of_applications": 1}}
	_, err = e.eventService.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": eventId},
		update,
	)

	return err
}

func (e *EventAppModel) GetEventsByUserEmail(email string) ([]*Event, error) {
	filter := bson.M{"attendee": bson.M{"$in": []string{email}}}

	eventApps, err := e.collection.Find(context.Background(), filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoRecords
		}
		return nil, err
	}

	var events []*Event
	for eventApps.Next(context.Background()) {
		var eventApp EventApps
		if err := eventApps.Decode(&eventApp); err != nil {
			return nil, err
		}

		eventObjID, err := primitive.ObjectIDFromHex(eventApp.EventID.Hex())
		if err != nil {
			return nil, err
		}

		event, err := e.eventService.GetEventByID(eventObjID)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				continue
			}
			return nil, err
		}
		events = append(events, event)
	}

	if err := eventApps.Err(); err != nil {
		return nil, err
	}
	return events, nil
}
