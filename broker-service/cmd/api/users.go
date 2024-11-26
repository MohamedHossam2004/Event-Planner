package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	request, err := http.NewRequest("POST", "http://authentication-service/v1/tokens/authentication", r.Body)
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

func (app *application) registerHandler(w http.ResponseWriter, r *http.Request) {
	request, err := http.NewRequest("POST", "http://authentication-service/v1/users", r.Body)
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

func (app *application) verifyTokenHandler(w http.ResponseWriter, r *http.Request) {
	request, err := http.NewRequest("POST", "http://authentication-service/v1/tokens/activation", r.Body)
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
