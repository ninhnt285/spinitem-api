package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"

	"../models"
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
	var shop models.ShopFull
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

// CallbackShop run callback for each platform
var CallbackShop http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
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
	// Run Callback Shop
	var shop *models.ShopFull
	switch platform {
	case "shopify":
		// Parse verify data from JSON
		var v models.ShopifyVerifier
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		// Get shop
		shop, err = models.GetShopBySession(v.State)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		// Run callback
		// shop.Shopify.Callback(v, shop)
		/*shop, err = v.Check()
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}*/
	}

	respondWithJSON(w, http.StatusOK, shop.Shop)
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
	shops, err := models.GetAllShops(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string][]models.Shop{"shops": shops})
}

// GetShop return a shop by id
var GetShop http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// Get userId
	userID, err := getUserID(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Can not get access")
		return
	}
	// Get shopID
	params := mux.Vars(r)
	shopID := params["id"]

	shop, err := models.GetShop(shopID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if shop.UserID != bson.ObjectId(userID) {
		respondWithError(w, http.StatusUnauthorized, "Can not get access")
		return
	}
	respondWithJSON(w, http.StatusOK, shop)
}
