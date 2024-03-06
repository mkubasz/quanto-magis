package main

import (
	"strings"
	"testing"
)

func TestRDDCreateFromArray(t *testing.T) {
	session := NewQuantoSession()
	data := []interface{}{1, 2, 3, 4, 5}
	rdd := session.Parallelize(data)
	if rdd == nil {
		t.Errorf("RDD is nil")
	}
	if rdd.size != 5 {
		t.Errorf("RDD size is not 5")
	}
}

func TestRDDCreateFromStringArray(t *testing.T) {
	session := NewQuantoSession()
	data := []interface{}{"A", "B", "C", "D", "E"}
	rdd := session.Parallelize(data)
	if rdd == nil {
		t.Errorf("RDD is nil")
	}
	if rdd.size != 5 {
		t.Errorf("RDD size is not 5")
	}
}

func lowerCase(s interface{}) interface{} {
	return strings.ToLower(s.(string))
}

func TestRDDShouldLowerCase(t *testing.T) {
	session := NewQuantoSession()
	data := []interface{}{"A", "B", "C", "D", "E"}
	rdd := session.Parallelize(data)
	result := rdd.Map(lowerCase)
	if result == nil {
		t.Errorf("result is nil")
	}
	if result.data[0] != "a" {
		t.Errorf("result[0] is not a")
	}
}

func TestRDDShouldFilterCLetter(t *testing.T) {
	session := NewQuantoSession()
	data := []interface{}{"A", "B", "C", "D", "E"}
	rdd := session.Parallelize(data)
	result := rdd.Filter(func(s interface{}) bool {
		return s.(string) != "C"
	})
	if result == nil {
		t.Errorf("result is nil")
	}
	if result.size != 4 {
		t.Errorf("result size is not 4")
	}
}

func TestRDDShouldFlatArray(t *testing.T) {
	session := NewQuantoSession()
	data := []interface{}{[]interface{}{1, 2, 3, 4, 5}, []interface{}{6, 7, 8, 9, 10}}

	rdd := session.Parallelize(data)
	result := rdd.FlatArray()
	if result == nil {
		t.Errorf("result is nil")
	}
	if result.size != 10 {
		t.Errorf("result size is not 10")
	}
}

func TestRDDShouldFlatMap(t *testing.T) {
	session := NewQuantoSession()
	data := []interface{}{[]interface{}{"A", "B", "C"}, []interface{}{"D", "E"}}

	rdd := session.Parallelize(data)
	result := rdd.FlatMap(lowerCase)
	if result == nil {
		t.Errorf("result is nil")
	}
	for _, d := range result.data {
		if d != "a" && d != "b" && d != "c" && d != "d" && d != "e" {
			t.Errorf("result data is not a, b, c, d, e")
		}
	}
}
