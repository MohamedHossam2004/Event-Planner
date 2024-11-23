package main

import (
	"encoding/json"
	"fmt"
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

	switch response.StatusCode {
	case http.StatusCreated:
		err = app.writeJSON(w, http.StatusCreated, payload, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	case http.StatusBadRequest:
		errMessage, ok := payload["error"].(string)
		if !ok {
			errMessage = "Invalid request"
		}
		app.badRequestResponse(w, r, fmt.Errorf(errMessage))
	case http.StatusUnprocessableEntity:
		errors, ok := payload["error"].(map[string]string)
		if !ok {
			app.serverErrorResponse(w, r, fmt.Errorf("invalid error format"))
			return
		}
		app.failedValidationResponse(w, r, errors)
	default:
		app.serverErrorResponse(w, r, fmt.Errorf("unexpected status code %d from authentication service", response.StatusCode))
	}
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

	switch response.StatusCode {
	case http.StatusCreated:
		err = app.writeJSON(w, http.StatusCreated, payload, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	case http.StatusBadRequest:
		errMessage, ok := payload["error"].(string)
		if !ok {
			errMessage = "Invalid request"
		}
		app.badRequestResponse(w, r, fmt.Errorf(errMessage))
	case http.StatusUnprocessableEntity:
		errors, ok := payload["error"].(map[string]string)
		if !ok {
			app.serverErrorResponse(w, r, fmt.Errorf("invalid error format"))
			return
		}
		app.failedValidationResponse(w, r, errors)
	default:
		app.serverErrorResponse(w, r, fmt.Errorf("unexpected status code %d from authentication service", response.StatusCode))
	}
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

	switch response.StatusCode {
	case http.StatusAccepted:
		err = app.writeJSON(w, http.StatusAccepted, payload, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	case http.StatusBadRequest:
		errMessage, ok := payload["error"].(string)
		if !ok {
			errMessage = "Invalid request"
		}
		app.badRequestResponse(w, r, fmt.Errorf(errMessage))
	case http.StatusUnprocessableEntity:
		errors, ok := payload["error"].(map[string]string)
		if !ok {
			app.serverErrorResponse(w, r, fmt.Errorf("invalid error format"))
			return
		}
		app.failedValidationResponse(w, r, errors)
	default:
		app.serverErrorResponse(w, r, fmt.Errorf("unexpected status code %d from authentication service", response.StatusCode))
	}
}
