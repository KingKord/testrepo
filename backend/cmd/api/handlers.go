package main

import (
	"fmt"
	"net/http"
)

func (app *Config) HeartBeat(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Successfully hit the backend",
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) PostNewUser(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Name       string `json:"name"`
		Surname    string `json:"surname"`
		Patronymic string `json:"patronymic"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Add to DB user %s %s %s", requestPayload.Name, requestPayload.Surname, requestPayload.Patronymic),
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
