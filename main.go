package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"./handlers"
	"./models"
)

func main() {
	models.ConnectDB("spinitem")
	defer models.CloseDB()

	router := mux.NewRouter()
	router.NotFoundHandler = handlers.NotFoundHandler
	router.Handle("/login", handlers.Login).Methods("POST")
	router.Handle("/register", handlers.Register).Methods("POST")
	// User handlers
	router.Handle("/users/{id}", handlers.Adapt(handlers.Auth, handlers.GetUser)).Methods("GET")
	// Item handlers
	router.Handle("/items", handlers.Adapt(handlers.Auth, handlers.AddItem)).Methods("POST")
	router.Handle("/items", handlers.Adapt(handlers.Auth, handlers.GetAllItems)).Methods("GET")
	router.Handle("/items/{id}", handlers.Adapt(handlers.Auth, handlers.GetItem)).Methods("GET")
	router.Handle("/items/{id}", handlers.Adapt(NotImplemented)).Methods("PUT")    // NotImplement
	router.Handle("/items/{id}", handlers.Adapt(NotImplemented)).Methods("DELETE") // NotImplement
	// Image handlers
	router.Handle("/images/upload", handlers.Adapt(handlers.Auth, handlers.UploadImage)).Methods("POST")
	router.Handle("/images", handlers.Adapt(handlers.Auth, handlers.AddImage)).Methods("POST")
	// Test handlers
	router.Handle("/test", NotImplemented).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}

// NotImplemented is dummy func for API Methods
var NotImplemented http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
}
