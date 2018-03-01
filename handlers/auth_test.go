package handlers

import (
	"fmt"
	"testing"

	"../models/user"
)

func TestGetToken(t *testing.T) {
	user := user.User{ID: "123456"}
	token, err := getToken(user)
	fmt.Println(token)
	if err != nil {
		t.Error(err.Error())
	}
}
