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

func TestShouldCreateBasicDataFrameWithTwoColumns(t *testing.T) {
	columnOne := []interface{}{"A", "B", "C", "D", "E"}
	columnTwo := []interface{}{1, 2, 3, 4, 5}

	df := NewDataFrame([]interface{}{columnOne, columnTwo}, []string{"col1", "col2"})
	if df == nil {
		t.Error("Failed to create DataFrame")
		return
	}
	if df.size != 10 {
		t.Error("Failed to create DataFrame")
	}
}
