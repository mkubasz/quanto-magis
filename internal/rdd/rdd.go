// Package rdd provides a Resilient Distributed Dataset (RDD) implementation
// for distributed data processing with parallel transformations.
package rdd

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
)

// RDD represents a Resilient Distributed Dataset, an immutable distributed
// collection of objects that can be processed in parallel.
// RDD uses Go generics for type safety.
type RDD[T any] struct {
	data []T
	size int
}

// New creates a new RDD from the provided data slice.
// The data is copied to ensure immutability.
func New[T any](data []T) *RDD[T] {
	// Create a copy to ensure immutability
	dataCopy := make([]T, len(data))
	copy(dataCopy, data)

	return &RDD[T]{
		data: dataCopy,
		size: len(dataCopy),
	}
}

// Collect returns all elements in the RDD as a slice.
// This is a terminal operation that returns a copy of the data.
func (r *RDD[T]) Collect() []T {
	result := make([]T, r.size)
	copy(result, r.data)
	return result
}

// Size returns the number of elements in the RDD.
func (r *RDD[T]) Size() int {
	return r.size
}

// Map applies the given function to each element in the RDD and returns
// a new RDD containing the transformed elements.
// Processing is done in parallel using worker goroutines.
//
// The context can be used to cancel the operation.
// Returns context.Canceled if the context is canceled during processing.
func (r *RDD[T]) Map(ctx context.Context, fn func(T) T) (*RDD[T], error) {
	// Check context before starting
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if r.size == 0 {
		return New([]T{}), nil
	}

	numWorkers := runtime.NumCPU()
	chunkSize := r.size / numWorkers
	if chunkSize == 0 {
		chunkSize = 1
		numWorkers = r.size
	}

	// Create channels for work distribution
	type job struct {
		index int
		value T
	}
	type result struct {
		index int
		value T
	}

	jobs := make(chan job, r.size)
	results := make(chan result, r.size)
	errors := make(chan error, numWorkers)

	// Start workers
	for i := 0; i < numWorkers; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					errors <- ctx.Err()
					return
				case j, ok := <-jobs:
					if !ok {
						return
					}
					// Apply transformation
					transformed := fn(j.value)
					select {
					case results <- result{index: j.index, value: transformed}:
					case <-ctx.Done():
						errors <- ctx.Err()
						return
					}
				}
			}
		}()
	}

	// Send jobs
	go func() {
		for i, val := range r.data {
			select {
			case jobs <- job{index: i, value: val}:
			case <-ctx.Done():
				close(jobs)
				return
			}
		}
		close(jobs)
	}()

	// Collect results
	mappedData := make([]T, r.size)
	for i := 0; i < r.size; i++ {
		select {
		case res := <-results:
			mappedData[res.index] = res.value
		case err := <-errors:
			return nil, err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return New(mappedData), nil
}

// Filter returns a new RDD containing only elements that satisfy the predicate.
// Processing is done in parallel using worker goroutines.
//
// The context can be used to cancel the operation.
// Returns context.Canceled if the context is canceled during processing.
func (r *RDD[T]) Filter(ctx context.Context, predicate func(T) bool) (*RDD[T], error) {
	// Check context before starting
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if r.size == 0 {
		return New([]T{}), nil
	}

	numWorkers := runtime.NumCPU()

	// Create channels for work distribution
	type job struct {
		value T
	}

	jobs := make(chan job, r.size)
	results := make(chan []T, numWorkers)
	errors := make(chan error, numWorkers)

	// Start workers
	for i := 0; i < numWorkers; i++ {
		go func() {
			var filtered []T
			for {
				select {
				case <-ctx.Done():
					errors <- ctx.Err()
					return
				case j, ok := <-jobs:
					if !ok {
						// Send local results to results channel
						results <- filtered
						return
					}
					// Apply predicate
					if predicate(j.value) {
						filtered = append(filtered, j.value)
					}
				}
			}
		}()
	}

	// Send jobs
	go func() {
		for _, val := range r.data {
			select {
			case jobs <- job{value: val}:
			case <-ctx.Done():
				close(jobs)
				return
			}
		}
		close(jobs)
	}()

	// Collect results from all workers
	var filteredData []T
	workersDone := 0
	for workersDone < numWorkers {
		select {
		case workerResults := <-results:
			filteredData = append(filteredData, workerResults...)
			workersDone++
		case err := <-errors:
			return nil, err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return New(filteredData), nil
}

// FlatArray flattens nested arrays/slices into a single-level RDD.
// For elements that are arrays or slices, it extracts all inner elements.
// For non-array elements, it includes them as-is.
//
// The context can be used to cancel the operation.
// Returns context.Canceled if the context is canceled during processing.
func (r *RDD[T]) FlatArray(ctx context.Context) (*RDD[T], error) {
	// Check context before starting
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if r.size == 0 {
		return New([]T{}), nil
	}

	numWorkers := runtime.NumCPU()
	chunkSize := r.size / numWorkers
	if chunkSize == 0 {
		chunkSize = 1
		numWorkers = r.size
	}

	results := make(chan []T, numWorkers)
	errors := make(chan error, numWorkers)

	// Start workers to process chunks
	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == numWorkers-1 {
			end = r.size
		}

		go func(start, end int) {
			var flattened []T

			for _, element := range r.data[start:end] {
				// Check context periodically
				select {
				case <-ctx.Done():
					errors <- ctx.Err()
					return
				default:
				}

				// Check if element is a slice or array using reflection
				t := reflect.TypeOf(element)
				if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
					v := reflect.ValueOf(element)
					// Extract all elements from the nested structure
					if innerSlice, ok := v.Interface().([]T); ok {
						flattened = append(flattened, innerSlice...)
					} else {
						// Fallback for type-incompatible slices
						flattened = append(flattened, element)
					}
				} else {
					flattened = append(flattened, element)
				}
			}

			results <- flattened
		}(start, end)
	}

	// Collect results from all workers
	var flattenedData []T
	workersDone := 0
	for workersDone < numWorkers {
		select {
		case workerResults := <-results:
			flattenedData = append(flattenedData, workerResults...)
			workersDone++
		case err := <-errors:
			return nil, err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return New(flattenedData), nil
}

// FlatMap flattens nested arrays and applies a transformation function to each element.
// It combines FlatArray and Map operations in a single pass for efficiency.
//
// The context can be used to cancel the operation.
// Returns context.Canceled if the context is canceled during processing.
func (r *RDD[T]) FlatMap(ctx context.Context, fn func(T) T) (*RDD[T], error) {
	// Check context before starting
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if r.size == 0 {
		return New([]T{}), nil
	}

	numWorkers := runtime.NumCPU()

	// Create channels for work distribution
	type job struct {
		value T
	}

	jobs := make(chan job, r.size)
	results := make(chan []T, numWorkers)
	errors := make(chan error, numWorkers)

	// Start workers
	for i := 0; i < numWorkers; i++ {
		go func() {
			var processed []T

			for {
				select {
				case <-ctx.Done():
					errors <- ctx.Err()
					return
				case j, ok := <-jobs:
					if !ok {
						// Send local results
						results <- processed
						return
					}

					// Check if element is a slice or array
					t := reflect.TypeOf(j.value)
					if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
						v := reflect.ValueOf(j.value)
						if innerSlice, ok := v.Interface().([]T); ok {
							// Process each element in the nested slice
							for _, inner := range innerSlice {
								transformed := fn(inner)
								processed = append(processed, transformed)
							}
						} else {
							// Fallback: treat as single element
							transformed := fn(j.value)
							processed = append(processed, transformed)
						}
					} else {
						// Process single element
						transformed := fn(j.value)
						processed = append(processed, transformed)
					}
				}
			}
		}()
	}

	// Send jobs
	go func() {
		for _, val := range r.data {
			select {
			case jobs <- job{value: val}:
			case <-ctx.Done():
				close(jobs)
				return
			}
		}
		close(jobs)
	}()

	// Collect results from all workers
	var flatMappedData []T
	workersDone := 0
	for workersDone < numWorkers {
		select {
		case workerResults := <-results:
			flatMappedData = append(flatMappedData, workerResults...)
			workersDone++
		case err := <-errors:
			return nil, err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return New(flatMappedData), nil
}

// String returns a string representation of the RDD.
func (r *RDD[T]) String() string {
	return fmt.Sprintf("RDD[size=%d]", r.size)
}
