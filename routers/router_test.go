package routers

import (
	"assignment/middleware"
	"bytes"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestNewRouter(t *testing.T) {
	testRouter := NewRouter()
	if testRouter == nil {
		t.Fatal("Router is nil")
	}
}

func TestHandleTransaction(t *testing.T) {
	err := godotenv.Load("./../.env")

	if err != nil {
		log.Fatalf("Error loading .env file: %+v", err)
	}
	middleware.URI = os.Getenv("MONGO_URL")

	var jsonStr = []byte(`{
    "user_id": "9999999aaaaaaaaaa",
    "currency": "EUR",
    "amount": 100.234,
    "time_placed": "24-JAN-21 10:27:44",
    "type": "deposit"
}`)
	req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	testRouter := NewRouter()
	testRouter.ServeHTTP(rr, req)

	stringResponse := rr.Body.String()
	if stringResponse != "Transaction stored!" {
		t.Errorf("Expected Transaction stored, got %s", stringResponse)
	}
}

func TestHandleBalance(t *testing.T) {
	err := godotenv.Load("./../.env")

	if err != nil {
		log.Fatalf("Error loading .env file: %+v", err)
	}
	middleware.URI = os.Getenv("MONGO_URL")

	req, _ := http.NewRequest("GET", "/balance?user_id=9999999aaaaaaaaaa", nil)

	rr := httptest.NewRecorder()
	testRouter := NewRouter()
	testRouter.ServeHTTP(rr, req)

	if rr.Body == nil {
		t.Errorf("Expected not nil response body")
	}
}
