package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/pascaldekloe/jwt"
)

func (app *application) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			app.authenticationRequiredResponse(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}
		token := headerParts[1]

		claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secret))

		if err != nil {
			app.notFoundResponse(w, r)
			return
		}

		if isAdmin, ok := claims.Set["isAdmin"].(bool); !ok || !isAdmin {
			app.notFoundResponse(w, r)
			return
		}

		if !claims.Valid(time.Now()) {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		if claims.Issuer != "giu-event-hub.com" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		if !claims.AcceptAudience("giu-event-hub.com") {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			app.authenticationRequiredResponse(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}
		token := headerParts[1]

		claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secret))
		if err != nil {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		if !claims.Valid(time.Now()) {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		if claims.Issuer != "giu-event-hub.com" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		if !claims.AcceptAudience("giu-event-hub.com") {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		if isActivated, ok := claims.Set["isActivated"].(bool); !ok || !isActivated {
			app.inactiveAccountResponse(w, r)
			return
		}

		if isAdmin, ok := claims.Set["isAdmin"].(bool); !ok || isAdmin {
			app.onlyUsersResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
