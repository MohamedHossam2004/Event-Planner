package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (h *application) routes() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Route("/v1/events", func(r chi.Router) {
		r.Get("/", h.getAllEventsHandler)                       // GET /events
		r.Get("/{id}", h.getEventByIDHandler)                   // GET /events/{id}
		r.Post("/", h.createEventHandler)                       // POST /events
		r.Put("/{id}", h.updateEventHandler)                    // PUT /events/{id}
		r.Delete("/{id}", h.deleteEventHandler)                 // DELETE /events/{id}
		r.Post("/{id}/apply", h.applyToEventHandler)            // POST /events/{id}/apply
		r.Delete("/{id}/unapply", h.removeUserEventApplication) // DELETE /events/{id}/unapply
	})

	r.Route("/v1/eventApps", func(r chi.Router) {
		r.Get("/", h.getAllEventAppsHandler)       // GET /eventApps
		r.Get("/{id}", h.getEventAppByIDHandler)   // GET /eventApps/{id}
		r.Post("/", h.createEventAppHandler)       // POST /eventApps
		r.Delete("/{id}", h.deleteEventAppHandler) // DELETE /eventApps/{id}
		r.Get("/user", h.viewAppliedEventsHandler) //GET /eventApps/user
	})

	return r
}
