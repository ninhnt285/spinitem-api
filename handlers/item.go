package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"

	"../models"
)

// AddItem handle for POST /items
var AddItem http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// Get userId
	userID, err := getUserID(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Can not get access")
		return
	}
	// Parse item from JSON
	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	item.UserID = bson.ObjectIdHex(userID)

	// Add new Item
	if err := item.Add(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, item)
}

// GetItem handle GET /items/{id}
var GetItem http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// Get itemId
	params := mux.Vars(r)
	itemID := params["id"]

	item, err := models.GetItem(itemID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, item)
}

// GetAllItems handle GET /items
var GetAllItems http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// Get userId
	userID, err := getUserID(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Can not get access")
		return
	}
	// Get all items
	items, err := models.GetAllItems(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string][]models.Item{"items": items})
}

// UpdateItem handle PUT /items/{id}
var UpdateItem http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var err error
	// Get userId
	userID, err := getUserID(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Can not get access")
		return
	}
	// Get itemId
	params := mux.Vars(r)
	itemID := params["id"]
	item, err := models.GetItem(itemID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	// Check owner permission
	if item.UserID != bson.ObjectIdHex(userID) {
		respondWithError(w, http.StatusUnauthorized, "Can not get access")
		return
	}
	// Parse item from JSON
	var updateItem models.Item
	if err := json.NewDecoder(r.Body).Decode(&updateItem); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	// Update item
	err = item.Update(updateItem)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, item)
}

// DeleteItem in DB
var DeleteItem http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var err error
	// Get userId
	userID, err := getUserID(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Can not get access")
		return
	}
	// Get Item
	params := mux.Vars(r)
	itemID := params["id"]
	item, err := models.GetItem(itemID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	// Check owner permission
	if item.UserID != bson.ObjectIdHex(userID) {
		respondWithError(w, http.StatusUnauthorized, "Can not get access")
		return
	}
	// Delete item
	err = item.Delete()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, item)
}
