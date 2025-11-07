package rdd

import (
	"reflect"
	"runtime"
	"sync"
)

type RDD[T any] struct {
	data []T
	size int
}

func New[T any](data []T) *RDD[T] {
	return &RDD[T]{
		data: data,
		size: len(data),
	}
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
	return New(mappedData)
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
	return New(filteredData)
}

func (r *RDD[T]) FlatArray() *RDD[T] {
	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()
	resultChan := make(chan T)

	// Fan-out
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(startIndex, endIndex int) {
			defer wg.Done()
			for _, d := range r.data[startIndex:endIndex] {
				if t := reflect.TypeOf(d); t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
					v := reflect.ValueOf(d)
					for _, inside := range v.Interface().([]T) {
						resultChan <- inside
					}
				} else {
					resultChan <- d
				}
			}
		}(i*len(r.data)/numWorkers, (i+1)*len(r.data)/numWorkers)
	}

	// Fan-in
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var flattenData []T
	for result := range resultChan {
		flattenData = append(flattenData, result)
	}

	return New(flattenData)
}

func (r *RDD[T]) Collect() []T {
	return r.data
}

func (r *RDD[T]) FlatMap(f func(T) T) *RDD[T] {
	numWorkers := runtime.NumCPU()
	jobs := make(chan T, len(r.data))
	results := make(chan []T, numWorkers)

	// Start worker goroutines
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			var flattenData []T
			for d := range jobs {
				t := reflect.TypeOf(d)
				if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
					v := reflect.ValueOf(d)
					insideSlice := v.Interface().([]T)
					for _, inside := range insideSlice {
						flattenData = append(flattenData, f(inside))
					}
				} else {
					flattenData = append(flattenData, f(d))
				}
			}
			results <- flattenData
		}()
	}

	go func() {
		for _, d := range r.data {
			jobs <- d
		}
		close(jobs)
	}()

	var flattenData []T
	go func() {
		for result := range results {
			flattenData = append(flattenData, result...)
		}
		close(results)
	}()

	wg.Wait()
	return New(flattenData)
}
