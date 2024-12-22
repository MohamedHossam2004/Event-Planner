package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/MohamedHossam2004/Event-Planner/event-service/internal/data"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (app *application) createEventAppHandler(w http.ResponseWriter, r *http.Request) {
	var eventApp struct {
		ID       primitive.ObjectID
		EventID  primitive.ObjectID
		Attendee []string
	}
	err := app.readJSON(w, r, &eventApp)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if eventApp.EventID.IsZero() {
		app.badRequestResponse(w, r, fmt.Errorf("EventID is required"))
		return
	}

	eventAppData := data.EventApps{
		ID:       eventApp.ID,
		EventID:  eventApp.EventID,
		Attendee: eventApp.Attendee,
	}

	err = app.models.EventApps.CreateEventApp(context.Background(), &eventAppData)
	if err != nil {
		app.Logger.Printf("Error creating event app: %v\n", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	event, err := app.models.Event.GetEventByID(eventApp.EventID)
	if err != nil {
		app.Logger.Printf("Error fetching event with ID %s: %v\n", eventApp.EventID.Hex(), err)
		app.badRequestResponse(w, r, err)
		return
	}

	payload := map[string]any{
		"event_name":     event.Name,
		"event_date":     event.Date,
		"event_location": fmt.Sprintf("%s,%s,%s,%s", event.Location.Address, event.Location.City, event.Location.State, event.Location.Country),
		"emails":         eventApp.Attendee,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		app.Logger.Printf("Error marshaling payload: %v\n", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.pushToQueue("event_register", string(jsonPayload))
	if err != nil {
		app.Logger.Printf("Error pushing to queue: %v\n", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"event_app": eventApp}, nil)
}

func (app *application) getAllEventAppsHandler(w http.ResponseWriter, r *http.Request) {
	eventApps, err := app.models.EventApps.ListEventApps(context.Background(), nil, nil)
	if err != nil {
		app.Logger.Printf("Error listing event apps: %v\n", err)
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "Failed to list event apps"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"event_apps": eventApps}, nil)
}

func (app *application) getEventAppByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "Invalid ID"}, nil)
		return
	}

	eventApp, err := app.models.EventApps.GetEventApp(context.Background(), objID)
	if err != nil {
		app.Logger.Printf("Error fetching event app with ID %s: %v\n", idStr, err)
		app.writeJSON(w, http.StatusNotFound, envelope{"error": "Event app not found"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"event_app": eventApp}, nil)
}

func (app *application) deleteEventAppHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "Invalid ID"}, nil)
		return
	}

	err = app.models.EventApps.DeleteEventApp(context.Background(), objID)
	if err != nil {
		app.Logger.Printf("Error deleting event app with ID %s: %v\n", idStr, err)
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "Failed to delete event app"}, nil)
		return
	}

	app.writeJSON(w, http.StatusNoContent, envelope{"message": "Event app deleted successfully"}, nil)
}

func (app *application) applyToEventHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "Invalid ID"}, nil)
		return
	}

	email, _, _, err := app.tokenExtractor.extractTokenData(r)
	if err != nil {
		app.writeJSON(w, http.StatusUnauthorized, envelope{"error": "Invalid token"}, nil)
		return
	}
	eventApp, err := app.models.EventApps.GetEventApp(context.Background(), objID)
	if err != nil {
		if errors.Is(err, data.ErrNoRecords) {
			app.writeJSON(w, http.StatusNotFound, envelope{"error": "Event app not found"}, nil)
			return
		}
		app.Logger.Printf("Error fetching event app with ID %s: %v\n", idStr, err)
		app.serverErrorResponse(w, r, err)
		return
	}

	if app.Contains(eventApp.Attendee, email) {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "You already applied to this event"}, nil)
		return
	}

	event, err := app.models.Event.GetEventByID(objID)
	if err != nil {
		if errors.Is(err, data.ErrNoRecords) {
			app.writeJSON(w, http.StatusNotFound, envelope{"error": "Event not found"}, nil)
			return
		}
		app.Logger.Printf("Error fetching event with ID %s: %v\n", objID.Hex(), err)
		app.serverErrorResponse(w, r, err)
		return
	}
	if event.Date.Before(time.Now()) {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "Event has ended"}, nil)
		return
	}

	err = app.models.EventApps.AddAttendeeToEvent(email, objID)
	if err != nil {
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "Failed to apply to event"}, nil)
		return
	}

	payload := map[string]any{
		"event_name":     event.Name,
		"event_date":     event.Date,
		"event_location": fmt.Sprintf("%s,%s,%s,%s", event.Location.Address, event.Location.City, event.Location.State, event.Location.Country),
		"emails":         []string{email},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		app.Logger.Printf("Error marshaling payload: %v\n", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.pushToQueue("event_register", string(jsonPayload))
	if err != nil {
		app.Logger.Printf("Error pushing to queue: %v\n", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "Applied to event successfully"}, nil)
}

func (app *application) removeUserEventApplication(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "Invalid ID"}, nil)
		return
	}

	email, _, _, err := app.tokenExtractor.extractTokenData(r)
	if err != nil {
		app.writeJSON(w, http.StatusUnauthorized, envelope{"error": "Invalid token"}, nil)
		return
	}

	eventApp, err := app.models.EventApps.GetEventApp(context.Background(), objID)
	if err != nil {
		if errors.Is(err, data.ErrNoRecords) {
			app.writeJSON(w, http.StatusNotFound, envelope{"error": "Event app not found"}, nil)
			return
		}
		app.Logger.Printf("Error fetching event app with ID %s: %v\n", idStr, err)
		app.serverErrorResponse(w, r, err)
		return
	}

	if !app.Contains(eventApp.Attendee, email) {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "You have not applied to this event"}, nil)
		return
	}

	event, err := app.models.Event.GetEventByID(objID)
	if err != nil {
		if errors.Is(err, data.ErrNoRecords) {
			app.writeJSON(w, http.StatusNotFound, envelope{"error": "Event not found"}, nil)
			return
		}
		app.Logger.Printf("Error fetching event with ID %s: %v\n", objID.Hex(), err)
		app.serverErrorResponse(w, r, err)
		return
	}
	if event.Date.Before(time.Now()) {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "Event has ended"}, nil)
		return
	}

	err = app.models.EventApps.RemoveAttendeeFromEvent(email, objID)
	if err != nil {
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "Failed to remove user event application"}, nil)
		return
	}
	app.writeJSON(w, http.StatusOK, envelope{"message": "Removed user event application successfully"}, nil)
}

func (app *application) viewAppliedEventsHandler(w http.ResponseWriter, r *http.Request) {
	email, _, _, err := app.tokenExtractor.extractTokenData(r)
	if err != nil {
		app.writeJSON(w, http.StatusUnauthorized, envelope{"error": "Invalid token"}, nil)
		return
	}

	events, err := app.models.EventApps.GetEventsByUserEmail(email)
	if err != nil {
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "Failed to fetch events"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"events": events}, nil)
}
