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
