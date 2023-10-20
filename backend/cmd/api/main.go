package main

import "net/http"

const webPort = 80

func main() {
	serv := &http.Server{}
	println("Starting web server on port", webPort)
	serv.ListenAndServe()
}
