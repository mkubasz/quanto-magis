package main

type Series[T any] struct {
	data []T
}

type DataFrame struct {
	series []Series[interface{}]
	size   int
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

func NewDataFrame(data []interface{}) *DataFrame {
	var series []Series[interface{}]
	size := 0
	for _, row := range data {
		var serie Series[interface{}]
		size += len(row.([]interface{}))
		serie.data = append(serie.data, row.([]interface{})...)
		series = append(series, serie)
	}
	return &DataFrame{
		size:   size,
		series: series,
	}
}
