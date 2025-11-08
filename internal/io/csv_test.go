package io_test

import (
	"testing"

	"mkubasz/quanto/internal/io"
)

func TestShouldReadCSVFile(t *testing.T) {
	reader := io.NewReader()
	df, err := reader.ReadCSV("../../testdata/test.csv")
	if err != nil {
		t.Errorf("failed to read csv file: %v", err)
	}
	expected := []string{"sepal.length", "sepal.width", "petal.length", "petal.width", "variety"}
	for _, e := range expected {
		if !df.HasColumn(e) {
			t.Errorf("expected %v in result data", e)
		}
	}
}
