package main

import (
	"fmt"
	"net/http"
)

//lint:ignore U1000 logError is used by error handling methods
func (app *application) logError(r *http.Request, err error) {
	app.Logger.Println(err)
}

//lint:ignore U1000 errorResponse is used by other error handling methods
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	res := envelope{"error": message}

	err := app.writeJSON(w, status, res, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

//lint:ignore U1000 serverErrorResponse is used by error handling methods
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

//lint:ignore U1000 notFoundResponse is used by error handling methods
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

//lint:ignore U1000 methodNotAllowedResponse is used by error handling methods
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

//lint:ignore U1000 editConflictResponse is used by error handling methods
func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, message)
}

//lint:ignore U1000 badRequestResponse is used by error handling methods
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

//lint:ignore U1000 failedValidationResponse is used by error handling methods
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

//lint:ignore U1000 invalidCredentialsResponse is used by error handling methods
func (app *application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

//lint:ignore U1000 invalidAuthenticationTokenResponse is used by error handling methods
func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

//lint:ignore U1000 authenticationRequiredResponse is used by error handling methods
func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message :=
		"you must be authenticated to access this resource"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

//lint:ignore U1000 inactiveAccountResponse is used by error handling methods
func (app *application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message :=
		"your user account must be activated to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}

//lint:ignore U1000 onlyAdminResponse is used by error handling methods
func (app *application) onlyAdminResponse(w http.ResponseWriter, r *http.Request) {
	message := "only administrators can access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}

//lint:ignore U1000 onlyUsersResponse is used by error handling methods
func (app *application) onlyUsersResponse(w http.ResponseWriter, r *http.Request) {
	message := "only authenticated users can access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}
