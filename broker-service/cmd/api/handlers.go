package main

import (
	"net/http"
)

func (app *application) Broker(w http.ResponseWriter, r *http.Request) {

	_ = app.writeJSON(w, http.StatusOK, envelope{"Broker": "Welcome to broker"}, nil)
}
