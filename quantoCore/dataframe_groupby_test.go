package main

import (
	"testing"
)

func TestShouldCountGroupBy(t *testing.T) {
	columnOne := []interface{}{"A", "B", "A", "D", "E"}
	columnTwo := []interface{}{1, 2, 3, 4, 5}
	df := NewDataFrame([]interface{}{columnOne, columnTwo}, []string{"col1", "col2"})
	dfg, _ := df.GroupBy("col1")
	new_df := dfg.Agg(Count).Show()
	if new_df.size != 2 {
		t.Errorf("Expected 2 columns, got %d", len(new_df.columns))
	}
}
