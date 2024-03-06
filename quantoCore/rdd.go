package main

import (
	"reflect"
	"runtime"
	"sync"
)

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

func (r *RDD[T]) AsyncTransform(fn func([]T, []T) []T) []T {
	var wg sync.WaitGroup
	numGoroutines := runtime.NumCPU()
	chunkSize := r.size / numGoroutines
	processedData := make([]T, 0, r.size)
	lock := sync.Mutex{}
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			start := i * chunkSize
			end := start + chunkSize
			if i == numGoroutines-1 {
				end = r.size
			}
			newData := r.data[start:end]
			lock.Lock()
			processedData = fn(processedData, newData)
			lock.Unlock()
		}(i)
	}
	wg.Wait()
	return processedData
}

func (r *RDD[T]) Map(f func(T) T) *RDD[T] {
	mappedData := r.AsyncTransform(func(data []T, newData []T) []T {
		for _, d := range newData {
			data = append(data, f(d))
		}
		return data
	})
	return RDDCreateFromArray(mappedData)
}

func (r *RDD[T]) Filter(f func(T) bool) *RDD[T] {
	filteredData := r.AsyncTransform(func(data []T, newData []T) []T {
		for _, d := range newData {
			if f(d) {
				data = append(data, d)
			}
		}
		return data
	})
	return RDDCreateFromArray(filteredData)
}

func (r *RDD[T]) FlatArray() *RDD[T] {
	var flattenData []T
	var wg sync.WaitGroup
	resultChan := make(chan T)
	for _, d := range r.data {
		wg.Add(1)
		go func(d T) {
			defer wg.Done()
			if t := reflect.TypeOf(d); t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
				v := reflect.ValueOf(d)
				for _, inside := range v.Interface().([]T) {
					resultChan <- inside
				}
			} else {
				resultChan <- d
			}
		}(d)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		flattenData = append(flattenData, result)
	}
	return RDDCreateFromArray(flattenData)
}

func (r *RDD[T]) Collect() []T {
	return r.data
}

func (r *RDD[T]) FlatMap(f func(T) T) *RDD[T] {
	var flattenData []T
	var wg sync.WaitGroup
	resultChan := make(chan T)
	for _, d := range r.data {
		wg.Add(1)
		go func(d T) {
			defer wg.Done()
			if t := reflect.TypeOf(d); t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
				v := reflect.ValueOf(d)
				for _, inside := range v.Interface().([]T) {
					resultChan <- f(inside)
				}
			} else {
				resultChan <- f(d)
			}
		}(d)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		flattenData = append(flattenData, result)
	}

	return RDDCreateFromArray(flattenData)
}
