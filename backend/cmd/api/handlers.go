package main

import (
	"backend/data"
	"encoding/json"
	"errors"
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

	user := data.User{
		Name:       requestPayload.Name,
		Surname:    requestPayload.Surname,
		Patronymic: requestPayload.Patronymic,
	}
	var jsonFromService struct {
		Count       int    `json:"count"`
		Name        string `json:"name"`
		Age         int    `json:"age,omitempty"`
		Gender      string `json:"gender,omitempty"`
		Probability int    `json:"probability,omitempty"`
	}
	fromOpenAPI, err := data.GetInfoFromOpenAPI(fmt.Sprintf("https://api.agify.io/?name=%s", user.Name))
	defer fromOpenAPI.Body.Close()

	err = json.NewDecoder(fromOpenAPI.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if fromOpenAPI.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("API returned non-OK status code: "+fromOpenAPI.Status), http.StatusBadRequest)
		return
	}
	user.Agify = jsonFromService.Age

	// in the end insert new user to the DB
	id, err := user.Insert(user)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Add to DB user %s %d %d", jsonFromService.Name, id),
		Data:    jsonFromService,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
