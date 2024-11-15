package main

import (
	"net/http"
)

func (app *application) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload, nil)
}
