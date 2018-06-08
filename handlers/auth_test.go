package handlers

import (
	"fmt"
	"testing"

	"../models"
)

func TestGetToken(t *testing.T) {
	user := models.User{ID: "123456"}
	token, err := getToken(user)
	fmt.Println(token)
	if err != nil {
		t.Error(err.Error())
	}
}
