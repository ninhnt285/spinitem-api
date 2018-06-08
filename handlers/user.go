package handlers

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"

	"../models"
)

// GetUser return user by id
// - If user_id == "me", return current user
var GetUser http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// Get userId
	params := mux.Vars(r)
	userID := params["id"]

	if userID == "me" {
		tokenUserID := context.Get(r, "user_id")
		if tokenUserID == nil {
			respondWithError(w, http.StatusUnauthorized, "Can not get access")
			return
		}
		userID = tokenUserID.(string)
	}

	user, err := models.GetUser(userID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}
