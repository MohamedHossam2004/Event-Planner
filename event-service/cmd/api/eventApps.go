package main

import (
	"context"
	"net/http"

	"github.com/MohamedHossam2004/Event-Planner/event-service/internal/data"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (app *application) createEventAppHandler(w http.ResponseWriter, r *http.Request) {
	var eventApp data.EventApps
	err := app.readJSON(w, r, &eventApp)
	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": err.Error()}, nil)
		return
	}

	err = app.EventApp.CreateEventApp(context.Background(), &eventApp)
	if err != nil {
		app.Logger.Printf("Error creating event app: %v\n", err)
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "Failed to create event app"}, nil)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"event_app": eventApp}, nil)
}

func (app *application) getAllEventAppsHandler(w http.ResponseWriter, r *http.Request) {
	eventApps, err := app.EventApp.ListEventApps(context.Background(), nil, nil)
	if err != nil {
		app.Logger.Printf("Error listing event apps: %v\n", err)
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "Failed to list event apps"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"event_apps": eventApps}, nil)
}

func (app *application) getEventAppByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "Invalid ID"}, nil)
		return
	}

	eventApp, err := app.EventApp.GetEventApp(context.Background(), objID)
	if err != nil {
		app.Logger.Printf("Error fetching event app with ID %s: %v\n", id, err)
		app.writeJSON(w, http.StatusNotFound, envelope{"error": "Event app not found"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"event_app": eventApp}, nil)
}

func (app *application) updateEventAppHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	objID, err := primitive.ObjectIDFromHex(id)
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

	err = app.EventApp.UpdateEventApp(context.Background(), objID, update)
	if err != nil {
		app.Logger.Printf("Error updating event app with ID %s: %v\n", id, err)
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "Failed to update event app"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "Event app updated successfully"}, nil)
}

func (app *application) deleteEventAppHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "Invalid ID"}, nil)
		return
	}

	err = app.EventApp.DeleteEventApp(context.Background(), objID)
	if err != nil {
		app.Logger.Printf("Error deleting event app with ID %s: %v\n", id, err)
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "Failed to delete event app"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "Event app deleted successfully"}, nil)
}
