package main

import (
	"backend/data"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
		Count       int     `json:"count"`
		Name        string  `json:"name"`
		Age         int     `json:"age,omitempty"`
		Gender      string  `json:"gender,omitempty"`
		Probability float64 `json:"probability,omitempty"`
		Country     []struct {
			CountryID   string  `json:"country_id"`
			Probability float32 `json:"probability"`
		} `json:"country,omitempty"`
	}

	// req to openAPI for age
	ageFromAPI, err := data.GetInfoFromOpenAPI(fmt.Sprintf("https://api.agify.io/?name=%s", user.Name))
	defer ageFromAPI.Body.Close()

	err = json.NewDecoder(ageFromAPI.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if ageFromAPI.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("API returned non-OK status code: "+ageFromAPI.Status), http.StatusBadRequest)
		return
	}
	user.Age = jsonFromService.Age
	log.Println("Successfully got age from open API")

	// req to openAPI for gender
	genderFromAPI, err := data.GetInfoFromOpenAPI(fmt.Sprintf("https://api.genderize.io/?name=%s", user.Name))
	defer genderFromAPI.Body.Close()

	err = json.NewDecoder(genderFromAPI.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if genderFromAPI.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("API returned non-OK status code: "+genderFromAPI.Status), http.StatusBadRequest)
		return
	}
	user.Gender = jsonFromService.Gender
	log.Println("Successfully got Gender from open API")
	// req to openAPI for nationality
	nationalityFromAPI, err := data.GetInfoFromOpenAPI(fmt.Sprintf("https://api.nationalize.io/?name=%s", user.Name))
	defer nationalityFromAPI.Body.Close()

	err = json.NewDecoder(nationalityFromAPI.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if nationalityFromAPI.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("API returned non-OK status code: "+nationalityFromAPI.Status), http.StatusBadRequest)
		return
	}

	var maxProbability float32
	maxProbability = 0.0
	mostProbableCountry := ""

	for _, country := range jsonFromService.Country {
		if country.Probability > maxProbability {
			maxProbability = country.Probability
			mostProbableCountry = country.CountryID
		}
	}
	user.Nationality = mostProbableCountry

	// in the end insert new user to the DB
	user.ID, err = user.Insert(user)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Add to DB user %s %d %s %s %d", jsonFromService.Name, user.Age, user.Gender, user.Nationality, user.ID),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	var user data.User
	users, err := user.GetAll()
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Get all users:",
		Data:    users,
	}
	app.writeJSON(w, http.StatusOK, payload)
}
