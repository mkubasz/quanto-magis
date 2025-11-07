package session

import (
	"testing"
)

func TestQuantoSessionBasicSettings(t *testing.T) {
	sess := New().
		SetAppName("Quanto Session").
		SetMode("local").
		GetOrCreate()

	if sess.AppName != "Quanto Session" {
		t.Errorf("AppName is not 'Quanto Session'")
	}
	if sess.Mode != Local {
		t.Errorf("Mode is not 'local'")
	}
}
