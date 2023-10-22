package main

import (
	"fmt"
	"net/http"
)

const webPort = "80"

type Config struct{}

func main() {
	app := Config{}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.getRoutes(),
	}

	println("Starting web server on port", webPort)
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
