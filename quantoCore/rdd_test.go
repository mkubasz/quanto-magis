package main

import (
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
