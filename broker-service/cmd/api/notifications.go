package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)



func (app *application) subscribe(w http.ResponseWriter, r *http.Request) {
	eventType := chi.URLParam(r, "eventType")
	if eventType == "" {
		app.badRequestResponse(w, r, errors.New("eventType"))
		return
	}
	
	request, err := http.NewRequest("POST", fmt.Sprintf("http://notification-service/subscribe/%s", eventType), r.Body)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	request.Header = r.Header

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var payload map[string]any
	if err := json.Unmarshal(responseBody, &payload); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.logger.Print(response.StatusCode)

	app.handleResponseStatus(w, r, response.StatusCode, payload)
}