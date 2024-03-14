package main

import (
	"testing"
)

func TestShouldCountGroupBy(t *testing.T) {
	columnOne := []interface{}{"A", "B", "A", "D", "E"}
	columnTwo := []interface{}{1, 2, 3, 4, 5}
	df := NewDataFrame([]interface{}{columnOne, columnTwo}, []string{"col1", "col2"})
	gb := df.GroupBy("col1").Agg(df.Count).Show()
}
