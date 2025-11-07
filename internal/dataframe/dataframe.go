package dataframe

import (
	"errors"

	"mkubasz/quanto/internal/rdd"
)

type Series[T any] struct {
	Data []T
}

type DataFrame struct {
	series  []Series[interface{}]
	columns []string
	size    int
}

func NewFromRDD[T any](r *rdd.RDD[T]) *DataFrame {
	var series Series[interface{}]
	for _, v := range r.Collect() {
		series.Data = append(series.Data, v)
	}
	return &DataFrame{
		size:   len(series.Data),
		series: []Series[interface{}]{series},
	}
}

func New(data []interface{}, columns []string) *DataFrame {
	var series []Series[interface{}]
	size := 0
	for _, row := range data {
		var serie Series[interface{}]
		size += len(row.([]interface{}))
		serie.Data = append(serie.Data, row.([]interface{})...)
		series = append(series, serie)
	}
	return &DataFrame{
		size:    size,
		series:  series,
		columns: columns,
	}
}

func (df *DataFrame) Select(columnName string) (Series[interface{}], error) {
	idx, err := df.getColumnIndex(columnName)
	if err != nil {
		return Series[interface{}]{}, errors.New("column not found")
	}
	return df.series[idx], nil
}

func (df *DataFrame) getColumnIndex(name string) (int, error) {
	for idx, column := range df.columns {
		if column == name {
			return idx, nil
		}
	}
	return -1, errors.New("column not found")
}

func (df *DataFrame) HasColumn(name string) bool {
	for _, column := range df.columns {
		if column == name {
			return true
		}
	}
	return false
}

func (s Series[T]) Distinct(key string) (Series[interface{}], error) {
	uniqueValues := make(map[interface{}]struct{})
	for _, serie := range s.Data {
		uniqueValues[serie] = struct{}{}
	}

	distinctValues := make([]interface{}, 0, len(uniqueValues))
	for k := range uniqueValues {
		distinctValues = append(distinctValues, k)
	}
	return Series[interface{}]{Data: distinctValues}, nil
}

func (s Series[T]) Count() int {
	return len(s.Data)
}
