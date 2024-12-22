package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/MohamedHossam2004/Event-Planner/event-service/rabbit"
	"github.com/pascaldekloe/jwt"
)

type envelope map[string]any

type payload struct {
	Topic string         `json:"topic"`
	Data  map[string]any `json:"data"`
}

type TokenExtractor interface {
	extractTokenData(r *http.Request) (string, bool, bool, error)
}

type realTokenExtractor struct {
	jwtSecret string
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {

		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.Contains(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) pushToQueue(name, msg string) error {
	emitter, err := rabbit.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	data := envelope{}

	err = json.Unmarshal([]byte(msg), &data)
	if err != nil {
		return err
	}

	payload := payload{
		Topic: name,
		Data:  data,
	}

	json, err := json.MarshalIndent(&payload, "", "\t")
	if err != nil {
		return err
	}

	err = emitter.Push(string(json), name)
	if err != nil {
		return err
	}

	app.Logger.Println("Pushed event to queue", name, string(json))
	return nil
}

func (app *application) background(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				app.Logger.Println(fmt.Errorf("%s", err))
			}
		}()
		fn()
	}()
}

func (rte *realTokenExtractor) extractTokenData(r *http.Request) (string, bool, bool, error) {
	token := r.Header.Get("Authorization")

	if token == "" {
		return "", false, false, errors.New("missing authorization header")
	}

	token = strings.TrimSpace(strings.Replace(token, "Bearer", "", 1))

	claims, err := jwt.HMACCheck([]byte(token), []byte(rte.jwtSecret))
	if err != nil {
		return "", false, false, err
	}

	if !claims.Valid(time.Now()) {
		return "", false, false, errors.New("invalid token")
	}

	userEmail, ok := claims.Set["email"].(string)
	if !ok {
		return "", false, false, errors.New("invalid token")
	}

	role, ok := claims.Set["isAdmin"].(bool)
	if !ok {
		return "", false, false, errors.New("invalid token")
	}

	isActivated, ok := claims.Set["isActivated"].(bool)
	if !ok {
		return "", false, false, errors.New("invalid token")
	}

	return userEmail, role, isActivated, nil
}

func (app *application) Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
