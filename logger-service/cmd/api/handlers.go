package main

import (
	"logger-service/data"
	"net/http"
	"tools"
)

type RequestPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	_ = app.ReadJSON(w, r, &requestPayload)

	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		_ = app.ErrorJSON(w, err)
		return
	}

	resp := tools.JsonResponse{
		Error:   false,
		Message: "logged",
	}

	_ = app.WriteJSON(w, http.StatusAccepted, resp)
}
