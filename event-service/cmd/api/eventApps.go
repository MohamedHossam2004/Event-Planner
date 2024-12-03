package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MohamedHossam2004/Event-Planner/event-service/internal/data"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (app *application) createEventAppHandler(w http.ResponseWriter, r *http.Request) {
	var eventApp data.EventApps
	err := app.readJSON(w, r, &eventApp)
	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": err.Error()}, nil)
		return
	}

	err = app.models.EventApps.CreateEventApp(context.Background(), &eventApp)
	if err != nil {
		app.Logger.Printf("Error creating event app: %v\n", err)
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "Failed to create event app"}, nil)
		return
	}

	event, err := app.models.Event.GetEventByID(eventApp.EventID)
	if err != nil {
		app.Logger.Printf("Error fetching event with ID %s: %v\n", eventApp.EventID.Hex(), err)
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "Failed to fetch event"}, nil)
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
	idStr := chi.URLParam(r, "id")
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

func (app *application) updateEventAppHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "Invalid ID"}, nil)
		return
	}

	var update map[string]interface{}
	err = app.readJSON(w, r, &update)
	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": err.Error()}, nil)
		return
	}

	err = app.models.EventApps.UpdateEventApp(context.Background(), objID, update)
	if err != nil {
		app.Logger.Printf("Error updating event app with ID %s: %v\n", idStr, err)
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "Failed to update event app"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "Event app updated successfully"}, nil)
}

func (app *application) deleteEventAppHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
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
