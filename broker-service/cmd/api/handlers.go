package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"tools"
)

const (
	authenticationServiceURL = "http://authentication-service/authenticate"
	loggerServiceURL         = "http://logger-service/log"
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
	payload := tools.JsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.WriteJSON(w, http.StatusOK, payload, nil)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var RequestPayload RequestPayload

	err := app.ReadJSON(w, r, &RequestPayload)
	if err != nil {
		_ = app.ErrorJSON(w, err)
		return
	}

	switch RequestPayload.Action {
	case "auth":
		app.authenticate(w, RequestPayload.Auth)
	case "log":
		app.logItem(w, RequestPayload.Log)
	default:
		_ = app.ErrorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, err := json.MarshalIndent(entry, "", "\t")
	if err != nil {
		_ = app.ErrorJSON(w, err)
		return
	}

	request, err := http.NewRequest(
		http.MethodPost,
		loggerServiceURL,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		_ = app.ErrorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		_ = app.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		_ = app.ErrorJSON(w, err)
		return
	}

	var payload tools.JsonResponse
	payload.Error = false
	payload.Message = "Log Created"

	_ = app.WriteJSON(w, http.StatusAccepted, payload)
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// TODO: change to just marshall after completion
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest(
		http.MethodPost,
		authenticationServiceURL,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		_ = app.ErrorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		_ = app.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		_ = app.ErrorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		_ = app.ErrorJSON(w, errors.New("error calling auth service"))
		return
	}

	var decodedResponse tools.JsonResponse
	err = json.NewDecoder(response.Body).Decode(&decodedResponse)
	if err != nil {
		_ = app.ErrorJSON(w, err)
		return
	}
	if decodedResponse.Error {
		_ = app.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	responsePayload := tools.JsonResponse{
		Error:   false,
		Message: "Authenticated!",
		Data:    decodedResponse.Data,
	}

	_ = app.WriteJSON(w, http.StatusAccepted, responsePayload)
}
