package rdd

import (
	"context"
	"strings"
	"testing"
)

// TestNew verifies RDD creation from various data types
func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		wantSize int
	}{
		{
			name:     "create from integer array",
			input:    []interface{}{1, 2, 3, 4, 5},
			wantSize: 5,
		},
		{
			name:     "create from string array",
			input:    []interface{}{"A", "B", "C", "D", "E"},
			wantSize: 5,
		},
		{
			name:     "create from empty array",
			input:    []interface{}{},
			wantSize: 0,
		},
		{
			name:     "create from single element",
			input:    []interface{}{42},
			wantSize: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rdd := New(tt.input)

			if rdd == nil {
				t.Error("expected non-nil RDD")
				return
			}

			got := len(rdd.Collect())
			if got != tt.wantSize {
				t.Errorf("size = %d, want %d", got, tt.wantSize)
			}
		})
	}
}

// TestMap verifies the Map transformation
func TestMap(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		mapper   func(interface{}) interface{}
		expected []interface{}
	}{
		{
			name:  "lowercase string transformation",
			input: []interface{}{"A", "B", "C"},
			mapper: func(s interface{}) interface{} {
				return strings.ToLower(s.(string))
			},
			expected: []interface{}{"a", "b", "c"},
		},
		{
			name:  "double integer transformation",
			input: []interface{}{1, 2, 3},
			mapper: func(n interface{}) interface{} {
				return n.(int) * 2
			},
			expected: []interface{}{2, 4, 6},
		},
		{
			name:  "identity transformation",
			input: []interface{}{1, 2, 3},
			mapper: func(n interface{}) interface{} {
				return n
			},
			expected: []interface{}{1, 2, 3},
		},
		{
			name:  "empty array",
			input: []interface{}{},
			mapper: func(n interface{}) interface{} {
				return n
			},
			expected: []interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			rdd := New(tt.input)
			result, err := rdd.Map(ctx, tt.mapper)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			got := result.Collect()
			if len(got) != len(tt.expected) {
				t.Errorf("length = %d, want %d", len(got), len(tt.expected))
				return
			}

			for i, val := range got {
				if val != tt.expected[i] {
					t.Errorf("element[%d] = %v, want %v", i, val, tt.expected[i])
				}
			}
		})
	}
}

// TestFilter verifies the Filter transformation
func TestFilter(t *testing.T) {
	tests := []struct {
		name      string
		input     []interface{}
		predicate func(interface{}) bool
		expected  []interface{}
	}{
		{
			name:  "filter out C letter",
			input: []interface{}{"A", "B", "C", "D", "E"},
			predicate: func(s interface{}) bool {
				return s.(string) != "C"
			},
			expected: []interface{}{"A", "B", "D", "E"},
		},
		{
			name:  "filter even numbers",
			input: []interface{}{1, 2, 3, 4, 5, 6},
			predicate: func(n interface{}) bool {
				return n.(int)%2 == 0
			},
			expected: []interface{}{2, 4, 6},
		},
		{
			name:  "filter none - all pass",
			input: []interface{}{1, 2, 3},
			predicate: func(n interface{}) bool {
				return true
			},
			expected: []interface{}{1, 2, 3},
		},
		{
			name:  "filter all - none pass",
			input: []interface{}{1, 2, 3},
			predicate: func(n interface{}) bool {
				return false
			},
			expected: []interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			rdd := New(tt.input)
			result, err := rdd.Filter(ctx, tt.predicate)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			got := result.Collect()
			if len(got) != len(tt.expected) {
				t.Errorf("length = %d, want %d", len(got), len(tt.expected))
				return
			}

			// Filter may not preserve order in parallel execution,
			// so check that all expected elements are present
			for _, expected := range tt.expected {
				found := false
				for _, actual := range got {
					if actual == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected element %v not found in result", expected)
				}
			}
		})
	}
}

// TestFlatArray verifies flattening of nested arrays
func TestFlatArray(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		wantSize int
	}{
		{
			name: "flatten two arrays",
			input: []interface{}{
				[]interface{}{1, 2, 3, 4, 5},
				[]interface{}{6, 7, 8, 9, 10},
			},
			wantSize: 10,
		},
		{
			name: "flatten mixed nested and non-nested",
			input: []interface{}{
				[]interface{}{1, 2},
				3,
				[]interface{}{4, 5},
			},
			wantSize: 5,
		},
		{
			name:     "empty array",
			input:    []interface{}{},
			wantSize: 0,
		},
		{
			name: "single nested array",
			input: []interface{}{
				[]interface{}{1, 2, 3},
			},
			wantSize: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			rdd := New(tt.input)
			result, err := rdd.FlatArray(ctx)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			got := len(result.Collect())
			if got != tt.wantSize {
				t.Errorf("size = %d, want %d", got, tt.wantSize)
			}
		})
	}
}

// TestFlatMap verifies the FlatMap transformation
func TestFlatMap(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		mapper   func(interface{}) interface{}
		expected []interface{}
	}{
		{
			name:  "flatten and lowercase",
			input: []interface{}{[]interface{}{"A", "B", "C"}, []interface{}{"D", "E"}},
			mapper: func(s interface{}) interface{} {
				return strings.ToLower(s.(string))
			},
			expected: []interface{}{"a", "b", "c", "d", "e"},
		},
		{
			name:  "flatten and double",
			input: []interface{}{[]interface{}{1, 2}, []interface{}{3, 4}},
			mapper: func(n interface{}) interface{} {
				return n.(int) * 2
			},
			expected: []interface{}{2, 4, 6, 8},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			rdd := New(tt.input)
			result, err := rdd.FlatMap(ctx, tt.mapper)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			got := result.Collect()
			if len(got) != len(tt.expected) {
				t.Errorf("length = %d, want %d", len(got), len(tt.expected))
				return
			}

			// FlatMap may not preserve order in parallel execution,
			// so check that all expected elements are present
			for _, expected := range tt.expected {
				found := false
				for _, actual := range got {
					if actual == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected element %v not found in result", expected)
				}
			}
		})
	}
}

// TestMapCancellation verifies that Map respects context cancellation
func TestMapCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Create large dataset to ensure cancellation is detected
	data := make([]interface{}, 10000)
	for i := range data {
		data[i] = i
	}

	rdd := New(data)
	_, err := rdd.Map(ctx, func(n interface{}) interface{} {
		return n.(int) * 2
	})

	if err == nil {
		t.Error("expected error due to context cancellation")
	}
	if err != context.Canceled {
		t.Errorf("expected context.Canceled error, got: %v", err)
	}
}

// TestFilterCancellation verifies that Filter respects context cancellation
func TestFilterCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Create large dataset
	data := make([]interface{}, 10000)
	for i := range data {
		data[i] = i
	}

	rdd := New(data)
	_, err := rdd.Filter(ctx, func(n interface{}) bool {
		return n.(int)%2 == 0
	})

	if err == nil {
		t.Error("expected error due to context cancellation")
	}
	if err != context.Canceled {
		t.Errorf("expected context.Canceled error, got: %v", err)
	}
}

// TestConcurrentOperations verifies race-free concurrent operations
func TestConcurrentOperations(t *testing.T) {
	// Run with: go test -race
	data := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	rdd := New(data)
	ctx := context.Background()

	// Run multiple operations concurrently
	done := make(chan bool, 3)

	go func() {
		_, err := rdd.Map(ctx, func(n interface{}) interface{} {
			return n.(int) * 2
		})
		if err != nil {
			t.Errorf("Map error: %v", err)
		}
		done <- true
	}()

	go func() {
		_, err := rdd.Filter(ctx, func(n interface{}) bool {
			return n.(int)%2 == 0
		})
		if err != nil {
			t.Errorf("Filter error: %v", err)
		}
		done <- true
	}()

	go func() {
		_, err := rdd.FlatArray(ctx)
		if err != nil {
			t.Errorf("FlatArray error: %v", err)
		}
		done <- true
	}()

	// Wait for all operations
	for i := 0; i < 3; i++ {
		<-done
	}
}

// BenchmarkMap benchmarks the Map operation
func BenchmarkMap(b *testing.B) {
	data := make([]interface{}, 10000)
	for i := range data {
		data[i] = i
	}
	rdd := New(data)
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = rdd.Map(ctx, func(n interface{}) interface{} {
			return n.(int) * 2
		})
	}
}

// BenchmarkFilter benchmarks the Filter operation
func BenchmarkFilter(b *testing.B) {
	data := make([]interface{}, 10000)
	for i := range data {
		data[i] = i
	}
	rdd := New(data)
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = rdd.Filter(ctx, func(n interface{}) bool {
			return n.(int)%2 == 0
		})
	}
}

// BenchmarkFlatMap benchmarks the FlatMap operation
func BenchmarkFlatMap(b *testing.B) {
	data := make([]interface{}, 100)
	for i := range data {
		data[i] = []interface{}{i, i + 1, i + 2}
	}
	rdd := New(data)
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = rdd.FlatMap(ctx, func(n interface{}) interface{} {
			return n.(int) * 2
		})
	}
}

// BenchmarkMapParallel benchmarks Map with parallel execution
func BenchmarkMapParallel(b *testing.B) {
	data := make([]interface{}, 10000)
	for i := range data {
		data[i] = i
	}
	rdd := New(data)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = rdd.Map(ctx, func(n interface{}) interface{} {
				return n.(int) * 2
			})
		}
	})
}
