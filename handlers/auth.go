package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"golang.org/x/crypto/bcrypt"

	"../helpers/config"
	"../helpers/validation"
	userModel "../models/user"
)

func getToken(user userModel.User) (string, error) {
	cf := config.GetInstance()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Unix() + 30*24*3600})
	tokenString, err := token.SignedString([]byte(cf.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Auth parses JWT token and save to Request context
var Auth http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	// Get token
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		respondWithError(w, http.StatusUnauthorized, "Can not get access")
		context.Set(r, "stopAdapt", true)
		return
	}
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	// Parse token to find JWTSecret
	conf := config.GetInstance()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(conf.JWTSecret), nil
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Can not get access")
		context.Set(r, "stopAdapt", true)
		return
	}
	// Get Claims data
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		respondWithError(w, http.StatusUnauthorized, "Can not get access")
		context.Set(r, "stopAdapt", true)
		return
	}
	// TODO: Validate exp

	// Save user_id to context
	context.Set(r, "user_id", claims["user_id"])
}

// Login returns a token or error
var Login http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// Parse user from JSON
	var user userModel.UserFull
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	// Validate email
	if !validation.ValidateEmail(user.Email) {
		respondWithError(w, http.StatusBadRequest, "Email is not valid input")
		return
	}
	// Validate password
	if !validation.ValidateStringInput(user.Password) {
		respondWithError(w, http.StatusBadRequest, "Password is not valid input")
		return
	}
	// Try to get user
	existedUser, err := userModel.GetByEmailOrUsername(user.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Email or password is invalid")
		return
	}
	// Match Bcrypt hash password
	err = bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(user.Password))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Email or password is invalid")
		return
	}
	// Get JWT Token
	token, err := getToken(existedUser.User)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

// Register return a token or error
var Register http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// Parse new user from JSON
	var user userModel.UserFull
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	// Validate username
	if !validation.ValidateStringInput(user.Username) {
		respondWithError(w, http.StatusBadRequest, "Username is not valid")
		return
	}
	// Validate email
	if !validation.ValidateEmail(user.Email) {
		respondWithError(w, http.StatusBadRequest, "Email is not valid input")
		return
	}
	// Validate password
	if !validation.ValidateStringInput(user.Password) {
		respondWithError(w, http.StatusBadRequest, "Password is not valid input")
		return
	}
	// Try to insert new user to DB
	err := user.Add()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	// Get JWT Token
	token, err := getToken(user.User)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}
