package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(60 * time.Second))
	mux.Use(middleware.Heartbeat("/ping"))
	mux.Use(middleware.Throttle(100))

	mux.Post("/", app.Broker)
	mux.Post("/v1/login", app.loginHandler)
	mux.Post("/v1/register", app.registerHandler)
	mux.Post("/v1/verify", app.verifyTokenHandler)

	mux.Group(func(r chi.Router) {
		r.Use(app.requireUser)
		r.Post("/v1/eventApps", app.createEventAppHandler)
		r.Put("/v1/eventApps/{id}", app.updateEventAppHandler)
		r.Delete("/v1/eventApps/{id}", app.deleteEventAppHandler)
	})

	mux.Get("/v1/events", app.getAllEventsHandler)
	mux.Get("/v1/events/{id}", app.getEventByIDHandler)

	mux.Group(func(r chi.Router) {
		r.Use(app.requireAdmin)
		r.Post("/v1/events", app.createEventHandler)
		r.Put("/v1/events/{id}", app.updateEventHandler)
		r.Delete("/v1/events/{id}", app.deleteEventHandler)
		r.Get("/v1/eventApps", app.getAllEventAppsHandler)
		r.Get("/v1/eventApps/{id}", app.getEventAppByIDHandler)
	})

	return mux
}
