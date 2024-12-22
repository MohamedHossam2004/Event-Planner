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

	mux.Get("/v1/events", app.getAllEventsHandler)
	mux.Get("/v1/events/{id}", app.getEventByIDHandler)
	mux.Post("/v1/events", app.createEventHandler)
	mux.Put("/v1/events/{id}", app.updateEventHandler)
	mux.Delete("/v1/events/{id}", app.deleteEventHandler)
	mux.Post("/v1/events/{id}/apply", app.applyToEventHandler)

	mux.Get("/v1/eventApps", app.getAllEventAppsHandler)
	mux.Get("/v1/eventApps/{id}", app.getEventAppByIDHandler)
	mux.Post("/v1/eventApps", app.createEventAppHandler)
	mux.Put("/v1/eventApps/{id}", app.updateEventAppHandler)
	mux.Delete("/v1/eventApps/{id}", app.deleteEventAppHandler)

	return mux
}
