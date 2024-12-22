package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MohamedHossam2004/Event-Planner/event-service/internal/data"
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
	idStr := r.PathValue("id")

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

	err = app.models.EventApps.CreateEventApp(context.Background(), &data.EventApps{
		ID:       primitive.NewObjectID(),
		EventID:  createdEvent.ID,
		Attendee: []string{},
	})
	if err != nil {
		app.Logger.Printf("Error creating event app: %v", err)
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "Failed to create event app"}, nil)
		return
	}

	location := fmt.Sprintf("%s,\n%s,\n%s,\n%s", createdEvent.Location.Address, createdEvent.Location.City, createdEvent.Location.State, createdEvent.Location.Country)

	payload := map[string]any{
		"event_type":        createdEvent.Type,
		"event_name":        createdEvent.Name,
		"event_date":        createdEvent.Date,
		"event_description": createdEvent.Description,
		"event_location":    location,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		app.Logger.Printf("Error marshaling payload: %v", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.pushToQueue("event_add", string(jsonPayload))
	if err != nil {
		app.Logger.Printf("Error pushing event to queue: %v", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"event": createdEvent}, nil)
}

func (app *application) updateEventHandler(w http.ResponseWriter, r *http.Request) {
	app.Logger.Println("UpdateEvent called")

	idStr := r.PathValue("id")
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

	eventApps, err := app.models.EventApps.GetEventApp(context.Background(), id)
	if err != nil {
		app.Logger.Printf("Error fetching event apps: %v", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	if eventApps != nil {
		emails := eventApps.Attendee

		payload := map[string]any{
			"emails":            emails,
			"event_name":        event.Name,
			"event_date":        event.Date,
			"event_description": event.Description,
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			app.Logger.Printf("Error marshaling payload: %v", err)
			app.serverErrorResponse(w, r, err)
			return
		}

		err = app.pushToQueue("event_update", string(jsonPayload))
		if err != nil {
			app.Logger.Printf("Error pushing event to queue: %v", err)
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	app.writeJSON(w, http.StatusOK, envelope{"event": updatedEvent}, nil)
}

func (app *application) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	app.Logger.Println("DeleteEvent called")
	idStr := r.PathValue("id")

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		app.Logger.Printf("Invalid ID format: %v", err)
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "Invalid ID format"}, nil)
		return
	}

	event, err := app.models.Event.GetEventByID(id)
	if err != nil {
		app.Logger.Printf("Error fetching event by ID: %v", err)
		app.writeJSON(w, http.StatusNotFound, envelope{"error": "Failed to fetch event"}, nil)
		return
	}
	if event != nil {
		err = app.models.Event.DeleteEvent(id)
		if err != nil {
			app.Logger.Printf("Error deleting event: %v", err)
			app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "Failed to delete event"}, nil)
			return
		}

		if event.Date.After(time.Now()) {
			eventApps, err := app.models.EventApps.GetEventApp(context.Background(), id)
			if err != nil {
				app.Logger.Printf("Error fetching event apps: %v", err)
				app.serverErrorResponse(w, r, err)
				return
			}

			emails := eventApps.Attendee

			payload := map[string]any{
				"emails":     emails,
				"event_name": event.Name,
				"event_date": event.Date,
			}

			jsonPayload, err := json.Marshal(payload)
			if err != nil {
				app.Logger.Printf("Error marshaling payload: %v", err)
				app.serverErrorResponse(w, r, err)
				return
			}

			err = app.pushToQueue("event_remove", string(jsonPayload))
			if err != nil {
				app.Logger.Printf("Error pushing event to queue: %v", err)
				app.serverErrorResponse(w, r, err)
				return
			}
		}
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "Event deleted successfully"}, nil)
}
