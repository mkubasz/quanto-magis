package dataframe

import (
	"testing"

	"mkubasz/quanto/internal/rdd"
)

func TestShouldConvertRDDToDataFrame(t *testing.T) {
	data := []interface{}{"A", "B", "C", "D", "E"}
	r := rdd.New(data)
	df := NewFromRDD(r)
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

	df := New([]interface{}{columnOne, columnTwo}, []string{"col1", "col2"})
	if df == nil {
		t.Error("Failed to create DataFrame")
		return
	}
	if expectedSize := 10; df.size != expectedSize {
		t.Errorf("Expected DataFrame size %d; got %d", expectedSize, df.size)
	}
}

func TestShouldSelectSelectedColumn(t *testing.T) {
	columnOne := []interface{}{"A", "B", "C", "D", "E"}
	columnTwo := []interface{}{1, 2, 3, 4, 5}
	df := New([]interface{}{columnOne, columnTwo}, []string{"col1", "col2"})
	col1, err := df.Select("col1")
	size := len(col1.Data)
	if size != 5 && err != nil {
		t.Error("Failed creating column")
	}
}

func TestShouldDistinctSelectedColumn(t *testing.T) {
	columnOne := []interface{}{"A", "B", "A", "D", "E"}
	columnTwo := []interface{}{1, 2, 3, 4, 5}
	df := New([]interface{}{columnOne, columnTwo}, []string{"col1", "col2"})
	col1, _ := df.Select("col1")
	distincted, err := col1.Distinct("A")
	size := len(distincted.Data)
	if size != 4 && err != nil {
		t.Error("Failed creating column")
	}
}

func TestShouldDistinctSelectedColumnAndCount(t *testing.T) {
	columnOne := []interface{}{"A", "B", "A", "D", "E"}
	columnTwo := []interface{}{1, 2, 3, 4, 5}
	df := New([]interface{}{columnOne, columnTwo}, []string{"col1", "col2"})
	col1, _ := df.Select("col1")
	distincted, err := col1.Distinct("A")
	size := distincted.Count()
	if size != 4 && err != nil {
		t.Error("Failed creating column")
	}
}
