package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func (h *application) routes() http.Handler {

mux := mux.NewRouter()


corsOptions := cors.New(cors.Options{
	AllowedOrigins:   []string{"*"}, 
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Content-Type", "Authorization"},
	AllowCredentials: true,
	MaxAge:           300, 
})

mux.HandleFunc("/events", h.getAllEventsHandler).Methods("GET")            
mux.HandleFunc("/events/{id}", h.getEventByIDHandler).Methods("GET")       
mux.HandleFunc("/events", h.createEventHandler).Methods("POST")            
mux.HandleFunc("/events/{id}", h.updateEventHandler).Methods("PUT")        
mux.HandleFunc("/events/{id}", h.deleteEventHandler).Methods("DELETE")     
mux.HandleFunc("/eventApps", h.getAllEventAppsHandler).Methods("GET")             // List all event apps
mux.HandleFunc("/eventApps/{id}", h.getEventAppByIDHandler).Methods("GET")        // Get an event app by ID
mux.HandleFunc("/eventApps", h.createEventAppHandler).Methods("POST")             // Create a new event app
mux.HandleFunc("/eventApps/{id}", h.updateEventAppHandler).Methods("PUT")         // Update an existing event app
mux.HandleFunc("/eventApps/{id}", h.deleteEventAppHandler).Methods("DELETE")      // Delete an event app by ID



return corsOptions.Handler(mux)
}