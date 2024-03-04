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

func (r *RDD[T]) Map(f func(T) T) *RDD[T] {
	var newData []T
	for _, d := range r.data {
		newData = append(newData, f(d))
	}
	return RDDCreateFromArray(newData)
}

func (r *RDD[T]) Filter(f func(T) bool) *RDD[T] {
	var newData []T
	for _, d := range r.data {
		if f(d) {
			newData = append(newData, d)
		}
	}
	return RDDCreateFromArray(newData)
}
