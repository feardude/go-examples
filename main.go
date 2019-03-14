package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Book represents a book
type Book struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author represents an author
type Author struct {
	Name string `json:"name"`
}

var idToBook = make(map[string]Book)

func main() {
	loadBooks()
	router := getRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	setContentType(w)
	json.NewEncoder(w).Encode(idToBook)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	setContentType(w)
	book := findBook(r)
	json.NewEncoder(w).Encode(book)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	setContentType(w)

	var book Book
	json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Int())
	idToBook[book.ID] = book

	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	setContentType(w)

	book := findBook(r)

	var updated Book
	json.NewDecoder(r.Body).Decode(&updated)
	updated.ID = book.ID
	idToBook[updated.ID] = updated

	json.NewEncoder(w).Encode(updated)
}

func loadBooks() {
	idToBook["1"] = Book{ID: "1", Title: "The dark tower", Author: &Author{Name: "Stephen King"}}
	idToBook["2"] = Book{ID: "2", Title: "11/22/63", Author: &Author{Name: "Stephen King"}}
}

func getRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books", createBook).Methods("POST")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	return router
}

func setContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func findBook(r *http.Request) *Book {
	params := mux.Vars(r)
	book := idToBook[params["id"]]
	return &book
}
