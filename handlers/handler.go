package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	"github.com/gorilla/context"
)

type returnError struct {
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

type returnData struct {
	Success bool         `json:"success"`
	Data    interface{}  `json:"data,omitempty"`
	Error   *returnError `json:"error,omitempty"`
}

// Adapt wraps all Adapter(s)
func Adapt(adapters ...http.Handler) http.Handler {
	var f http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		for _, adapter := range adapters {
			if context.Get(r, "stopAdapt") != nil {
				r.Body.Close()
				break
			}
			adapter.ServeHTTP(w, r)
		}
	}
	return f
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, returnData{Success: false, Error: &returnError{Message: msg, Code: code}})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	var result returnData
	if reflect.TypeOf(payload).Name() == "returnData" {
		result = payload.(returnData)
	} else {
		result = returnData{Success: true, Data: payload}
	}

	response, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func getUserID(r *http.Request) (string, error) {
	// Get userId
	userIDObject := context.Get(r, "user_id")
	if userIDObject == nil {
		return "", errors.New("Can not get access")
	}
	userID, ok := userIDObject.(string)
	if !ok {
		return "", errors.New("Can not get access")
	}
	return userID, nil
}

// NotImplemented is dummy func for API Methods
var NotImplemented http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
}

// NotFoundHandler handle 404 request
var NotFoundHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusBadRequest, "Bad Request")
}
