package main

import "errors"

type Series[T any] struct {
	data []T
}

type DataFrame struct {
	series  []Series[interface{}]
	columns []string
	size    int
}

func (s *RDD[T]) ToDF() *DataFrame {
	var series Series[interface{}]
	for _, v := range s.data {
		series.data = append(series.data, v)
	}
	return &DataFrame{
		size:   s.size,
		series: []Series[interface{}]{series},
	}
}

func NewDataFrame(data []interface{}, columns []string) *DataFrame {
	var series []Series[interface{}]
	size := 0
	for _, row := range data {
		var serie Series[interface{}]
		size += len(row.([]interface{}))
		serie.data = append(serie.data, row.([]interface{})...)
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
		return Series[interface{}]{}, err
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

func (s Series[T]) Distinct(key string) (Series[interface{}], error) {
	uniqueValues := make(map[interface{}]struct{})
	for _, serie := range s.data {
		uniqueValues[serie] = struct{}{}
	}

	distinctValues := make([]interface{}, 0, len(uniqueValues))
	for k := range uniqueValues {
		distinctValues = append(distinctValues, k)
	}
	return Series[interface{}]{data: distinctValues}, nil
}

func (s Series[T]) Count() int {
	return len(s.data)
}
