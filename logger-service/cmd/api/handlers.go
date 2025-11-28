package main

import (
	"logger-service/data"
	"net/http"
)

type JsonPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var payload JsonPayload
	err := app.readJSON(w, r, &payload)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	event := data.LogEntry{
		Name: payload.Name,
		Data: payload.Data,
	}

	err = app.Models.LogEntry.Create(event)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	okPayload := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.writeJSON(w, http.StatusAccepted, okPayload)
}
