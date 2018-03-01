package handlers

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	itemModel "../models/item"
	"github.com/gorilla/mux"
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
	var item itemModel.Item
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

	item, err := itemModel.GetByID(itemID)
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
	items, err := itemModel.GetAll(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, items)
}
