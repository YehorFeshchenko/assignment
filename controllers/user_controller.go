package controllers

import (
	"assignment/middleware"
	"assignment/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func HandleTransaction(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: Handle Transaction")

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var transactionMsg models.Transaction
	err = json.Unmarshal(b, &transactionMsg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	isSuccessful := middleware.StoreTransaction(transactionMsg)
	if !isSuccessful {
		http.Error(w, "Failed to store data", 500)
		return
	}

	w.Header().Set("content-type", "text")
	_, err = w.Write([]byte("Transaction stored!"))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func HandleBalance(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: Handle Balance")

	id := r.URL.Query().Get("user_id")
	balance := middleware.GetUserBalance(id)

	balanceJSON, err := json.Marshal(balance)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("content-type", "text")
	_, err = w.Write(balanceJSON)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
