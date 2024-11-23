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

	mux.NotFound(http.HandlerFunc(app.notFoundResponse))
	mux.MethodNotAllowed(http.HandlerFunc(app.methodNotAllowedResponse))

	mux.Post("/v1/users", app.registerUserHandler)
	mux.Post("/v1/tokens/activation", app.createActivationTokenHandler)
	mux.Put("/v1/users/activated", app.activateUserHandler)
	mux.Post("/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return mux
}
