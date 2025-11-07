package rdd

import (
	"strings"
	"testing"
)

func TestRDDCreateFromArray(t *testing.T) {
	data := []interface{}{1, 2, 3, 4, 5}
	r := New(data)
	if r == nil {
		t.Errorf("RDD is nil")
		return
	}
	if len(r.Collect()) != 5 {
		t.Errorf("RDD size is not 5")
	}
}

func TestRDDCreateFromStringArray(t *testing.T) {
	data := []interface{}{"A", "B", "C", "D", "E"}
	r := New(data)
	if r == nil {
		t.Errorf("RDD is nil")
		return
	}
	if len(r.Collect()) != 5 {
		t.Errorf("RDD size is not 5")
	}
}

func lowerCase(s interface{}) interface{} {
	return strings.ToLower(s.(string))
}

func TestRDDShouldLowerCase(t *testing.T) {
	data := []interface{}{"A", "B", "C", "D", "E"}
	r := New(data)
	expected := []string{"a", "b", "c", "d", "e"}
	if len(r.Collect()) != len(expected) {
		t.Errorf("expected result size to be %d, got %d", len(expected), len(r.Collect()))
	}
}

func TestRDDShouldFilterCLetter(t *testing.T) {
	data := []interface{}{"A", "B", "C", "D", "E"}
	r := New(data)
	result := r.Filter(func(s interface{}) bool {
		return s.(string) != "C"
	})
	if len(result.Collect()) != 4 {
		t.Errorf("expected result size to be 4, got %d", len(result.Collect()))
	}
}

func TestRDDShouldFlatArray(t *testing.T) {
	data := []interface{}{
		[]interface{}{1, 2, 3, 4, 5},
		[]interface{}{6, 7, 8, 9, 10},
	}

	r := New(data)
	result := r.FlatArray()

	if len(result.Collect()) != 10 {
		t.Errorf("expected result size to be 10, got %d", len(result.Collect()))
	}
}

func TestRDDShouldFlatMap(t *testing.T) {
	data := []interface{}{[]interface{}{"A", "B", "C"}, []interface{}{"D", "E"}}

	r := New(data)
	result := r.FlatMap(lowerCase)

	if result == nil {
		t.Errorf("result is nil")
		return
	}
	expected := []interface{}{"a", "b", "c", "d", "e"}
	resultData := result.Collect()
	for _, e := range expected {
		if !contains(resultData, e) {
			t.Errorf("expected %v in result data", e)
		}
	}
}

func contains(slice []interface{}, val interface{}) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
