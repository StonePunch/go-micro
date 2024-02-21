package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"tools"
)

const loggerServiceURL = "http://logger-service/log"

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.ReadJSON(w, r, &requestPayload)
	if err != nil {
		_ = app.ErrorJSON(w, err)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		_ = app.ErrorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		_ = app.ErrorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	// create log when user is successfully authenticated
	err = app.logRequest("User Authenticated", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		_ = app.ErrorJSON(w, errors.New("failed to log user authentication"), http.StatusInternalServerError)
	}

	payload := tools.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("user '%s' logged in", user.Email),
		Data:    user,
	}

	_ = app.WriteJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logRequest(name, data string) error {
	var logEntry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	logEntry.Name = name
	logEntry.Data = data

	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", loggerServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
