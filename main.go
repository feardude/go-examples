package main

import (
	"encoding/json"
	"log"
	"net/http"

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

var books []Book

func main() {
	loadBooks()
	router := getRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	setContentType(w)
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	setContentType(w)
	params := mux.Vars(r)
	for _, book := range books {
		if book.ID == params["id"] {
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

func loadBooks() {
	books = append(books, Book{ID: "1", Title: "The dark tower", Author: &Author{Name: "Stephen King"}})
	books = append(books, Book{ID: "2", Title: "11/22/63", Author: &Author{Name: "Stephen King"}})
}

func getRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	return router
}

func setContentType(w http.ResponseWriter) {
	// writer := *w
	w.Header().Set("Content-Type", "application/json")
}
