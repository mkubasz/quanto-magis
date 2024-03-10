package main

import (
	"testing"
)

func TestShouldReadCSVFile(t *testing.T) {
	session := NewQuantoSession().
	SetAppName("Quanto Session").
	SetMode("local").
	GetOrCreate()
	session.Read.csv("data/iris.csv")
}
