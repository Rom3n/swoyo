package main

import (
	"net/http"
	"swoyo/controllers"
	"swoyo/utils"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	utils.DatabaseClient = utils.ConnectDB()
	router.HandleFunc("/", utils.SetCorsHeaders(controllers.EncodeUrls)).Methods("POST")
	router.HandleFunc("/", utils.SetCorsHeaders(controllers.DecodeUrl)).Methods("GET")
	http.ListenAndServe(":8080", router)
}
