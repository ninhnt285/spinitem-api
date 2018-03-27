package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"

	shopify "../helpers/shopify"
	shopModel "../models/shop"
)

// AddShop handle for POST /shops
var AddShop http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// Get userId
	userID, err := getUserID(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Can not get access")
		return
	}
	// Parse Shop from JSON
	var shop shopModel.ShopFull
	if err := json.NewDecoder(r.Body).Decode(&shop); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	shop.UserID = bson.ObjectIdHex(userID)

	// Add new Shop
	if err := shop.Add(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, shop)
}

// VerifyShop save AccessCode and send some test APIs
var VerifyShop http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// Get platform
	params := mux.Vars(r)
	platform := params["platform"]
	// Get userId
	_, err := getUserID(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Can not get access")
		return
	}

	// Verify Shop
	var newShop *shopModel.ShopFull
	switch platform {
	case "shopify":
		// Parse verify data from JSON
		var v shopify.Verify
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		newShop, err = v.Check()
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	respondWithJSON(w, http.StatusOK, newShop.Shop)
}

// GetAllShops return all shops of the user
var GetAllShops http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// Get userId
	userID, err := getUserID(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Can not get access")
		return
	}
	// Get all shops
	shops, err := shopModel.GetAll(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string][]shopModel.Shop{"shops": shops})
}
