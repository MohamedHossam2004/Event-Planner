package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()


	mux.HandleFunc("GET /v1/events/", app.getAllEventsHandler)                       // GET /events
	mux.HandleFunc("GET /v1/events/{id}", app.getEventByIDHandler)                   // GET /events/{id}
	mux.HandleFunc("POST /v1/events", app.createEventHandler)                       // POST /events
	mux.HandleFunc("PUT /v1/events/{id}", app.updateEventHandler)                    // PUT /events/{id}
	mux.HandleFunc("DELETE /v1/events/{id}", app.deleteEventHandler)                 // DELETE /events/{id}
	mux.HandleFunc("POST /v1/events/{id}/apply", app.applyToEventHandler)            // POST /events/{id}/apply
	mux.HandleFunc("DELETE /v1/events/{id}/unapply", app.removeUserEventApplication) // DELETE /events/{id}/unapply

	mux.HandleFunc("GET /v1/eventApps/", app.getAllEventAppsHandler)       // GET /eventApps
	mux.HandleFunc("GET /v1/eventApps/{id}", app.getEventAppByIDHandler)   // GET /eventApps/{id}
	mux.HandleFunc("POST /v1/eventApps", app.createEventAppHandler)       // POST /eventApps
	mux.HandleFunc("DELETE /v1/eventApps/{id}", app.deleteEventAppHandler) // DELETE /eventApps/{id}
	mux.HandleFunc("GET /v1/eventApps/user", app.viewAppliedEventsHandler) //GET /eventApps/user

	// Routes
	// r.Route("/v1/events", func(r chi.Router) {
	// 	r.Get("/", h.getAllEventsHandler)                       // GET /events
	// 	r.Get("/{id}", h.getEventByIDHandler)                   // GET /events/{id}
	// 	r.Post("/", h.createEventHandler)                       // POST /events
	// 	r.Put("/{id}", h.updateEventHandler)                    // PUT /events/{id}
	// 	r.Delete("/{id}", h.deleteEventHandler)                 // DELETE /events/{id}
	// 	r.Post("/{id}/apply", h.applyToEventHandler)            // POST /events/{id}/apply
	// 	r.Delete("/{id}/unapply", h.removeUserEventApplication) // DELETE /events/{id}/unapply
	// })

	// r.Route("/v1/eventApps", func(r chi.Router) {
	// 	r.Get("/", h.getAllEventAppsHandler)       // GET /eventApps
	// 	r.Get("/{id}", h.getEventAppByIDHandler)   // GET /eventApps/{id}
	// 	r.Post("/", h.createEventAppHandler)       // POST /eventApps
	// 	r.Delete("/{id}", h.deleteEventAppHandler) // DELETE /eventApps/{id}
	// 	r.Get("/user", h.viewAppliedEventsHandler) //GET /eventApps/user
	// })

	return logger(recoverer(cors(mux)))
}
