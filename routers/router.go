package routers

import (
	"assignment/controllers"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	httpRouter := mux.NewRouter().StrictSlash(true)
	httpRouter.HandleFunc("/transactions", controllers.HandleTransaction)
	httpRouter.HandleFunc("/balance", controllers.HandleBalance)
	return httpRouter
}
