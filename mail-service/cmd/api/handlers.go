package main

import (
	"fmt"
	"net/http"
	"tools"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	err := app.ReadJSON(w, r, &requestPayload)
	if err != nil {
		_ = app.ErrorJSON(w, err)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		_ = app.ErrorJSON(w, err)
		return
	}

	responsePayload := tools.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Email send to %s successfully!", requestPayload.To),
	}

	_ = app.WriteJSON(w, http.StatusAccepted, responsePayload)
}
