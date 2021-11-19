package main

import (
	"log"
	"net/http"

	"assignment/routers"
)

func main() {
	router := routers.NewRouter()
	log.Fatal(http.ListenAndServe(":4000", router))
}
