package main

import "net/http"

func (app *Config) HeartBeat(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Successfully hit the backend",
	}

	app.writeJSON(w, http.StatusOK, payload)
}
