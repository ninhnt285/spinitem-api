package main

import (
	"log"
	"net/http"

	goHandlers "github.com/gorilla/handlers"

	"github.com/gorilla/mux"

	"./handlers"
	"./models"
)

func main() {
	models.ConnectDB("spinitem")
	defer models.CloseDB()

	router := mux.NewRouter()
	router.NotFoundHandler = handlers.NotFoundHandler
	// File handlers
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./public/"))))
	// Auth handlers
	router.Handle("/login", handlers.Login).Methods("POST")
	router.Handle("/register", handlers.Register).Methods("POST")
	// User handlers
	router.Handle("/users/{id}", handlers.Adapt(handlers.Auth, handlers.GetUser)).Methods("GET")
	// Shop handlers
	router.Handle("/shops", handlers.Adapt(handlers.Auth, handlers.GetAllShops)).Methods("GET")
	router.Handle("/shops", handlers.Adapt(handlers.Auth, handlers.AddShop)).Methods("POST")
	router.Handle("/shops/verify/{platform}", handlers.Adapt(handlers.Auth, handlers.VerifyShop)).Methods("POST")
	// Item handlers
	router.Handle("/items", handlers.Adapt(handlers.Auth, handlers.AddItem)).Methods("POST")
	router.Handle("/items", handlers.Adapt(handlers.Auth, handlers.GetAllItems)).Methods("GET")
	router.Handle("/items/{id}", handlers.GetItem).Methods("GET")
	router.Handle("/items/{id}", handlers.Adapt(handlers.Auth, handlers.UpdateItem)).Methods("PUT")
	router.Handle("/items/{id}", handlers.Adapt(handlers.Auth, handlers.DeleteItem)).Methods("DELETE")
	// Image handlers
	router.Handle("/images", handlers.GetAllImages).Methods("GET")
	router.Handle("/images/upload", handlers.Adapt(handlers.Auth, handlers.UploadImage)).Methods("POST")
	router.Handle("/images", handlers.Adapt(handlers.Auth, handlers.AddImage)).Methods("POST")
	// Test handlers
	router.Handle("/test", handlers.NotImplemented).Methods("GET")

	headersOk := goHandlers.AllowedHeaders([]string{"Authorization", "Content-Type"})
	originsOk := goHandlers.AllowedOrigins([]string{"*"})
	methodsOk := goHandlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	log.Fatal(http.ListenAndServe(":8000", goHandlers.CORS(originsOk, headersOk, methodsOk)(router)))
}
