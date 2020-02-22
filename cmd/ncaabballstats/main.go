package main

import (
	"log"
	"ncaabballstats/internal/handlers"
	"net/http"
)

func main() {
	a := &handlers.App{
		TeamHandler: new(handlers.TeamHandler),
		ApiHandler:  new(handlers.ApiHandler),
	}

	log.Fatal(http.ListenAndServe(":8080", a))
}
