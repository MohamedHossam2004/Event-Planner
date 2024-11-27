package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"errors"
	"github.com/go-chi/chi/v5"
)

func (app *application) getAllEventsHandler(w http.ResponseWriter, r *http.Request) {
	request, err := http.NewRequest("GET", "http://event-service/v1/events", nil)
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

	app.handleResponseStatus(w, r, response.StatusCode, payload)
}

func (app *application) getEventByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		app.badRequestResponse(w, r, errors.New("missing id"))
		return
	}

	request, err := http.NewRequest("GET", fmt.Sprintf("http://event-service/v1/events/%s", idStr), nil)
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

	app.handleResponseStatus(w, r, response.StatusCode, payload)
}

func (app *application) createEventHandler(w http.ResponseWriter, r *http.Request) {
	request, err := http.NewRequest("POST", "http://event-service/v1/events", r.Body)
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

	app.handleResponseStatus(w, r, response.StatusCode, payload)
}

func (app *application) updateEventHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		app.badRequestResponse(w, r, errors.New("missing id"))
		return
	}

	request, err := http.NewRequest("PUT", fmt.Sprintf("http://event-service/v1/events/%s", idStr), r.Body)
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

	app.handleResponseStatus(w, r, response.StatusCode, payload)
}

func (app *application) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		app.badRequestResponse(w, r, errors.New("missing id"))
		return
	}

	request, err := http.NewRequest("DELETE", fmt.Sprintf("http://event-service/v1/events/%s", idStr), nil)
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

	app.handleResponseStatus(w, r, response.StatusCode, payload)
}

func (app *application) getAllEventAppsHandler(w http.ResponseWriter, r *http.Request) {
	request, err := http.NewRequest("GET", "http://event-service/v1/eventApps", nil)
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

	app.handleResponseStatus(w, r, response.StatusCode, payload)
}

func (app *application) getEventAppByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		app.badRequestResponse(w, r, errors.New("missing id"))
		return
	}

	request, err := http.NewRequest("GET", fmt.Sprintf("http://event-service/v1/eventApps/%s", idStr), nil)
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

	app.handleResponseStatus(w, r, response.StatusCode, payload)
}

func (app *application) createEventAppHandler(w http.ResponseWriter, r *http.Request) {
	request, err := http.NewRequest("POST", "http://event-service/v1/eventApps", r.Body)
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

	app.handleResponseStatus(w, r, response.StatusCode, payload)
}

func (app *application) updateEventAppHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		app.badRequestResponse(w, r, errors.New("missing id"))
		return
	}

	request, err := http.NewRequest("PUT", fmt.Sprintf("http://event-service/v1/eventApps/%s", idStr), r.Body)
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

	app.handleResponseStatus(w, r, response.StatusCode, payload)
}

func (app *application) deleteEventAppHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		app.badRequestResponse(w, r, errors.New("missing id"))
		return
	}

	request, err := http.NewRequest("DELETE", fmt.Sprintf("http://event-service/v1/eventApps/%s", idStr), nil)
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

	app.handleResponseStatus(w, r, response.StatusCode, payload)
}
