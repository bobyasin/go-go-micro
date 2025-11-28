package main

import (
	"log"
	"net/http"
)

const (
	publishPort = "80"
)

type Config struct{}

func main() {
	app := Config{}

	log.Printf("Starting broker service on port %s\n", publishPort)

	srv := &http.Server{
		Addr:    ":" + publishPort,
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}

}
