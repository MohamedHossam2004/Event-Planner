package main

import (
	"net/http"

	"github.com/MohamedHossam2004/Event-Planner/event-service/internal/data"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (app *application) getAllEventsHandler(w http.ResponseWriter, r *http.Request) {
	app.Logger.Println("GetAllEvents called")

	events, err := app.models.Event.GetAllEvents()
	if err != nil {
		app.Logger.Printf("Error fetching events: %v", err)
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "Failed to fetch events"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"events": events}, nil)
}

func (app *application) getEventByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		app.Logger.Printf("Invalid ID format: %v", err)
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "Invalid ID format"}, nil)
		return
	}

	event, err := app.models.Event.GetEventByID(id)
	if err != nil {
		app.Logger.Printf("Error fetching event by ID: %v", err)
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "Failed to fetch event"}, nil)
		return
	}
	if event == nil {
		app.Logger.Println("Event not found")
		app.writeJSON(w, http.StatusNotFound, envelope{"error": "Event not found"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"event": event}, nil)
}

func (app *application) createEventHandler(w http.ResponseWriter, r *http.Request) {
	app.Logger.Println("CreateEvent called")

	var event data.Event
	if err := app.readJSON(w, r, &event); err != nil {
		app.Logger.Printf("Error decoding event: %v", err)
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": err.Error()}, nil)
		return
	}

	createdEvent, err := app.models.Event.CreateEvent(&event)
	if err != nil {
		app.Logger.Printf("Error creating event: %v", err)
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "Failed to create event"}, nil)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"event": createdEvent}, nil)
}

func (app *application) updateEventHandler(w http.ResponseWriter, r *http.Request) {
	app.Logger.Println("UpdateEvent called")

	idStr := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		app.Logger.Printf("Invalid ID format: %v", err)
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "Invalid ID format"}, nil)
		return
	}

	var event data.Event
	if err := app.readJSON(w, r, &event); err != nil {
		app.Logger.Printf("Error decoding event: %v", err)
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": err.Error()}, nil)
		return
	}

	updatedEvent, err := app.models.Event.UpdateEvent(id, &event)
	if err != nil {
		app.Logger.Printf("Error updating event: %v", err)
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "Failed to update event"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"event": updatedEvent}, nil)
}

func (app *application) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	app.Logger.Println("DeleteEvent called")
	idStr := chi.URLParam(r, "id")

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		app.Logger.Printf("Invalid ID format: %v", err)
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "Invalid ID format"}, nil)
		return
	}

	err = app.models.Event.DeleteEvent(id)
	if err != nil {
		app.Logger.Printf("Error deleting event: %v", err)
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "Failed to delete event"}, nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
