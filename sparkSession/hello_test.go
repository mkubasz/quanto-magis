package main

import (
	"testing"
)

func TestQuantumSessionBasicSettings(t *testing.T) {
	session := NewQuantumSession().
		SetAppName("Quantum Session").
		master("local").
		getOrCreate()

	if session.AppName != "Quantum Session" {
		t.Errorf("AppName is not 'Quantum Session'")
	}
	if session.Mode != Local {
		t.Errorf("Mode is not 'local'")
	}
}
