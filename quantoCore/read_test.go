package main

import (
	"testing"
)

func TestShouldReadCSVFile(t *testing.T) {
	session := NewQuantoSession().
		SetAppName("Quanto Session").
		SetMode("local").
		GetOrCreate()
	df := session.Read.Csv("data/iris.csv")
	expected := []string{"sepal.length", "sepal.width", "petal.length", "petal.width", "variety"}
	for _, e := range expected {
		if !containsStr(df.columns, e) {
			t.Errorf("expected %v in result data", e)
		}
	}
}
