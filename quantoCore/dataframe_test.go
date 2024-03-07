package main

import (
	"testing"
)

func TestShouldConvertRDDToDataFrame(t *testing.T) {
	session := NewQuantoSession().GetOrCreate()
	data := []interface{}{"A", "B", "C", "D", "E"}
	rdd := session.Parallelize(data)
	df := rdd.ToDF()
	if df == nil {
		t.Error("Failed to convert RDD to DataFrame")
		return
	}
	if df.size != 5 {
		t.Error("Failed to convert RDD to DataFrame")
	}
}
