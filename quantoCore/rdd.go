package main

type RDD[T any] struct {
	data []T
	size int
}

func RDDCreateFromArray[T any](data []T) *RDD[T] {
	return &RDD[T]{
		data: data,
		size: len(data),
	}
}

func (s *QuantoSession) Parallelize(data []interface{}) *RDD[interface{}] {
	return RDDCreateFromArray(data)
}
