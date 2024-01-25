package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload, nil)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var RequestPayload RequestPayload

	err := app.readJSON(w, r, &RequestPayload)
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	switch RequestPayload.Action {
	case "auth":
		app.authenticate(w, RequestPayload.Auth)
	case "log":
		app.logItem(w, RequestPayload.Log)
	default:
		_ = app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, err := json.MarshalIndent(entry, "", "\t")
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	logServiceURL := "http://logger-service/log"
	request, err := http.NewRequest(
		"POST",
		logServiceURL,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		_ = app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	_ = app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create payload for the auth service
	authRequestPayload, _ := json.MarshalIndent(a, "", "\t")

	// call auth service to validate user
	authRequest, err := http.NewRequest(
		"POST",
		"http://authentication-service/authenticate",
		bytes.NewBuffer(authRequestPayload),
	)
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	authResponse, err := client.Do(authRequest)
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}
	defer authResponse.Body.Close()

	// validate response from auth service
	if authResponse.StatusCode == http.StatusUnauthorized {
		_ = app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if authResponse.StatusCode != http.StatusAccepted {
		_ = app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	// read response data
	var decodedResponse jsonResponse
	err = json.NewDecoder(authResponse.Body).Decode(&decodedResponse)
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	if decodedResponse.Error {
		_ = app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	// respond to initial authentication request
	responsePayload := jsonResponse{
		Error:   false,
		Message: "Authenticated!",
		Data:    decodedResponse.Data,
	}

	_ = app.writeJSON(w, http.StatusAccepted, responsePayload)
}
