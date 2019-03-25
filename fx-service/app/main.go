package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Started FX service")
	InitDB()

	defer ShutdownDB()

	router := getRouter()
	log.Fatal(http.ListenAndServe(":8080", router))

	log.Println("Finished FX service")
}

func getCurrencies(w http.ResponseWriter, r *http.Request) {
	setContentType(w)
	currencies := GetCurrencies()
	json.NewEncoder(w).Encode(currencies)
}

func getRate(w http.ResponseWriter, r *http.Request) {

}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/currencies", getCurrencies).Methods("GET")
	router.HandleFunc("/currencies/{code}", getRate).Methods("GET")
	return router
}

func setContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
