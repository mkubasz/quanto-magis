package main

import (
	"testing"
)

func TestQuantoSessionBasicSettings(t *testing.T) {
	session := NewQuantoSession().
		SetAppName("Quanto Session").
		master("local").
		getOrCreate()

	if session.AppName != "Quanto Session" {
		t.Errorf("AppName is not 'Quanto Session'")
	}
	if session.Mode != Local {
		t.Errorf("Mode is not 'local'")
	}
}
