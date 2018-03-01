package config

import (
	"testing"
)

func TestGetInstance(t *testing.T) {
	c := GetInstance()
	if c == nil {
		t.Error("Can not get config data")
	}

	d := GetInstance()
	if d.JWTSecret != "nMVyPedy2JooPnlg85vQ" {
		t.Errorf("%s is not equal %s", d.JWTSecret, "nMVyPedy2JooPnlg85vQ")
	}
}
