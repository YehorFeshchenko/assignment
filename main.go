package main

import (
	"assignment/middleware"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"assignment/routers"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	middleware.URI = os.Getenv("MONGO_URL")

	router := routers.NewRouter()
	log.Println("Server listening on port :4000")
	log.Fatal(http.ListenAndServe(":4000", router))
}
