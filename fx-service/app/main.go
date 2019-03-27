package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Started FX service")
	defer log.Println("Finished FX service")

	// TODO: Graceful shutdown
	InitDB()
	defer ShutdownDB()

	router := getRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getCurrencies(w http.ResponseWriter, r *http.Request) {
	currencies := GetCurrencies()
	jsonResponse(w, currencies)
}

func getRate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	base := strings.ToUpper(params["base"])
	rate := GetRate(base, time.Now())
	rate.Quote = resolveQuoteCode(params)
	jsonResponse(w, rate)
}

func resolveQuoteCode(params map[string]string) string {
	quote := params["quote"]
	if quote == "" {
		quote = "RUB"
	}
	return strings.ToUpper(quote)
}

func jsonResponse(w http.ResponseWriter, object interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(object)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/currencies", getCurrencies).Methods("GET")
	router.HandleFunc("/currencies/{base}", getRate).Methods("GET")
	return router
}
