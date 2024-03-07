package main

import (
	"strings"
	"testing"
)

func TestRDDCreateFromArray(t *testing.T) {
	session := NewQuantoSession().GetOrCreate()
	data := []interface{}{1, 2, 3, 4, 5}
	rdd := session.Parallelize(data)
	if rdd == nil {
		t.Errorf("RDD is nil")
		return
	}
	if rdd.size != 5 {
		t.Errorf("RDD size is not 5")
	}
}

func TestRDDCreateFromStringArray(t *testing.T) {
	session := NewQuantoSession().GetOrCreate()
	data := []interface{}{"A", "B", "C", "D", "E"}
	rdd := session.Parallelize(data)
	if rdd == nil {
		t.Errorf("RDD is nil")
		return
	}
	if rdd.size != 5 {
		t.Errorf("RDD size is not 5")
	}
}

func lowerCase(s interface{}) interface{} {
	return strings.ToLower(s.(string))
}

func TestRDDShouldLowerCase(t *testing.T) {
	session := NewQuantoSession().GetOrCreate()
	data := []interface{}{"A", "B", "C", "D", "E"}
	rdd := session.Parallelize(data)
	expected := []string{"a", "b", "c", "d", "e"}
	if rdd.size != len(expected) {
		t.Errorf("expected result size to be %d, got %d", len(expected), rdd.size)
	}
}

func TestRDDShouldFilterCLetter(t *testing.T) {
	session := NewQuantoSession().GetOrCreate()
	data := []interface{}{"A", "B", "C", "D", "E"}
	rdd := session.Parallelize(data)
	result := rdd.Filter(func(s interface{}) bool {
		return s.(string) != "C"
	})
	if result.size != 4 {
		t.Errorf("expected result size to be 4, got %d", result.size)
	}
}

func TestRDDShouldFlatArray(t *testing.T) {
	session := NewQuantoSession().GetOrCreate()
	data := []interface{}{
		[]interface{}{1, 2, 3, 4, 5},
		[]interface{}{6, 7, 8, 9, 10},
	}

	rdd := session.Parallelize(data)
	result := rdd.FlatArray()

	if len(result.data) != 10 {
		t.Errorf("expected result size to be 10, got %d", len(result.data))
	}
}

func TestRDDShouldFlatMap(t *testing.T) {
	session := NewQuantoSession().GetOrCreate()
	data := []interface{}{[]interface{}{"A", "B", "C"}, []interface{}{"D", "E"}}

	rdd := session.Parallelize(data)
	result := rdd.FlatMap(lowerCase)

	if result == nil {
		t.Errorf("result is nil")
		return
	}
	expected := []interface{}{"a", "b", "c", "d", "e"}
	for _, e := range expected {
		if !contains(result.data, e) {
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
