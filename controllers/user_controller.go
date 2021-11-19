package controllers

import (
	"assignment/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func HandleTransaction(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println(transactionMsg)

	transactionJSON, err := json.Marshal(transactionMsg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(transactionJSON)

	fmt.Println("Endpoint Hit: Handle Transaction")
}
