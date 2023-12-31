package main

import (
	"backend/data"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (app *Config) HeartBeat(w http.ResponseWriter, r *http.Request) {
	log.Println("hit the / (home) page")
	payload := jsonResponse{
		Error:   false,
		Message: "Successfully hit the backend",
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) PostNewUser(w http.ResponseWriter, r *http.Request) {
	log.Println("hit the /new page")
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
	log.Println("hit the /users page")
	// Get parameter "gender" from URL
	gender := r.URL.Query().Get("gender")
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	var user data.User
	var users []*data.User
	var err error

	if gender != "" {
		if page != "" && limit != "" {
			p, _ := strconv.Atoi(page)
			l, _ := strconv.Atoi(limit)
			users, err = user.GetAllUsersByGender(gender, l, p)
		} else if page != "" {
			p, _ := strconv.Atoi(page)
			users, err = user.GetAllUsersByGender(gender, 2, p)
		} else if limit != "" {
			l, _ := strconv.Atoi(limit)
			users, err = user.GetAllUsersByGender(gender, l, 1)
		} else {
			users, err = user.GetAllUsersByGender(gender, 100, 1)
		}
	} else {
		if page != "" && limit != "" {
			p, _ := strconv.Atoi(page)
			l, _ := strconv.Atoi(limit)
			users, err = user.GetAll(l, p)
		} else if page != "" {
			p, _ := strconv.Atoi(page)
			users, err = user.GetAll(2, p)
		} else if limit != "" {
			l, _ := strconv.Atoi(limit)
			users, err = user.GetAll(l, 1)
		} else {
			users, err = user.GetAll(100, 1)
		}
	}

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	if len(users) == 0 {
		app.errorJSON(w, errors.New("no records in DB"))
	}
	log.Printf("Successfully got %d records from DB", len(users))
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Total Users: %d", len(users)),
		Data:    users,
	}
	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) DeleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete URL hit")

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var user data.User

	err = user.DeleteByID(id)
	log.Printf("Deleted user in DB under id %d\n", id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("User under ID:%d deleted", id),
	}
	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) UpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Update URL hit")
	var requestPayload struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Surname     string `json:"surname"`
		Patronymic  string `json:"patronymic,omitempty"`
		Age         string `json:"age,omitempty"`
		Gender      string `json:"gender,omitempty"`
		Nationality string `json:"nationality,omitempty"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	id, _ := strconv.Atoi(requestPayload.ID)
	age, _ := strconv.Atoi(requestPayload.Age)
	user := data.User{
		ID:          id,
		Name:        requestPayload.Name,
		Surname:     requestPayload.Surname,
		Patronymic:  requestPayload.Patronymic,
		Age:         age,
		Gender:      requestPayload.Age,
		Nationality: requestPayload.Nationality,
	}
	err = user.Update()
	log.Printf("Updated user under id %d", user.ID)
	if err != nil {
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("User under ID:%d updated", user.ID),
	}
	app.writeJSON(w, http.StatusOK, payload)
}
