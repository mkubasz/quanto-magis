package dataframe

import (
	"context"
	"errors"
	"testing"
)

// TestGroupBy verifies grouping operations
func TestGroupBy(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() *DataFrame
		columnName string
		wantErr    error
		wantGroups int
	}{
		{
			name: "group by existing column",
			setup: func() *DataFrame {
				df, _ := New(
					[]interface{}{
						[]interface{}{"A", "B", "A", "C", "B"},
						[]interface{}{1, 2, 3, 4, 5},
					},
					[]string{"col1", "col2"},
				)
				return df
			},
			columnName: "col1",
			wantErr:    nil,
			wantGroups: 3, // A, B, C
		},
		{
			name: "group by non-existing column",
			setup: func() *DataFrame {
				df, _ := New(
					[]interface{}{
						[]interface{}{1, 2, 3},
					},
					[]string{"col1"},
				)
				return df
			},
			columnName: "col2",
			wantErr:    ErrColumnNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			df := tt.setup()
			grouped, err := df.GroupBy(ctx, tt.columnName)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(grouped.groups) != tt.wantGroups {
				t.Errorf("groups = %d, want %d", len(grouped.groups), tt.wantGroups)
			}
		})
	}
}

// TestGroupByAgg verifies aggregation operations
func TestGroupByAgg(t *testing.T) {
	ctx := context.Background()

	df, _ := New(
		[]interface{}{
			[]interface{}{"A", "B", "A", "C", "B", "A"},
			[]interface{}{1, 2, 3, 4, 5, 6},
		},
		[]string{"category", "value"},
	)

	grouped, err := df.GroupBy(ctx, "category")
	if err != nil {
		t.Fatalf("GroupBy failed: %v", err)
	}

	// Test Count aggregation
	result, err := grouped.Agg(Count).Show(ctx)
	if err != nil {
		t.Fatalf("Show failed: %v", err)
	}

	if result == nil {
		t.Error("expected non-nil result")
	}
}

// TestGroupByMultipleAgg verifies multiple aggregations
func TestGroupByMultipleAgg(t *testing.T) {
	ctx := context.Background()

	df, _ := New(
		[]interface{}{
			[]interface{}{"A", "B", "A"},
			[]interface{}{1, 2, 3},
		},
		[]string{"category", "value"},
	)

	grouped, err := df.GroupBy(ctx, "category")
	if err != nil {
		t.Fatalf("GroupBy failed: %v", err)
	}

	// Add multiple aggregations
	result, err := grouped.Agg(Count).Agg(Sum).Show(ctx)
	if err != nil {
		t.Fatalf("Show failed: %v", err)
	}

	if result == nil {
		t.Error("expected non-nil result")
	}
}

// TestGroupByShowWithoutAgg verifies error when Show is called without aggregations
func TestGroupByShowWithoutAgg(t *testing.T) {
	ctx := context.Background()

	df, _ := New(
		[]interface{}{
			[]interface{}{"A", "B", "A"},
		},
		[]string{"category"},
	)

	grouped, err := df.GroupBy(ctx, "category")
	if err != nil {
		t.Fatalf("GroupBy failed: %v", err)
	}

	// Try to Show without adding aggregations
	_, err = grouped.Show(ctx)
	if err == nil {
		t.Error("expected error when Show is called without aggregations")
	}
	if !errors.Is(err, ErrInvalidData) {
		t.Errorf("expected ErrInvalidData, got %v", err)
	}
}

// TestGroupByCancellation verifies context cancellation
func TestGroupByCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Create large dataset
	data := make([]interface{}, 10000)
	for i := range data {
		data[i] = i % 100
	}

	df, _ := New(
		[]interface{}{data},
		[]string{"col1"},
	)

	_, err := df.GroupBy(ctx, "col1")
	if err == nil {
		t.Error("expected error due to context cancellation")
	}
	if err != context.Canceled {
		t.Errorf("expected context.Canceled error, got: %v", err)
	}
}

// TestShowCancellation verifies context cancellation during Show
func TestShowCancellation(t *testing.T) {
	ctx := context.Background()

	// Create large dataset
	data := make([]interface{}, 10000)
	for i := range data {
		data[i] = i % 100
	}

	df, _ := New(
		[]interface{}{data},
		[]string{"col1"},
	)

	grouped, err := df.GroupBy(ctx, "col1")
	if err != nil {
		t.Fatalf("GroupBy failed: %v", err)
	}

	// Cancel context before Show
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = grouped.Agg(Count).Show(ctx)
	if err == nil {
		t.Error("expected error due to context cancellation")
	}
	if err != context.Canceled {
		t.Errorf("expected context.Canceled error, got: %v", err)
	}
}

// TestAggregationFunctions verifies built-in aggregation functions
func TestAggregationFunctions(t *testing.T) {
	tests := []struct {
		name     string
		fn       func([]interface{}) int
		input    []interface{}
		expected int
	}{
		{
			name:     "count function",
			fn:       Count,
			input:    []interface{}{1, 2, 3, 4, 5},
			expected: 5,
		},
		{
			name:     "count empty",
			fn:       Count,
			input:    []interface{}{},
			expected: 0,
		},
		{
			name:     "sum function",
			fn:       Sum,
			input:    []interface{}{1, 2, 3, 4, 5},
			expected: 15,
		},
		{
			name:     "sum with non-integers",
			fn:       Sum,
			input:    []interface{}{1, "string", 2, 3},
			expected: 6, // Non-integers ignored
		},
		{
			name:     "sum empty",
			fn:       Sum,
			input:    []interface{}{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn(tt.input)
			if result != tt.expected {
				t.Errorf("result = %d, want %d", result, tt.expected)
			}
		})
	}
}

// BenchmarkGroupBy benchmarks grouping operation
func BenchmarkGroupBy(b *testing.B) {
	ctx := context.Background()

	// Create dataset with 1000 rows and 10 unique groups
	data := make([]interface{}, 1000)
	for i := range data {
		data[i] = i % 10
	}

	df, _ := New(
		[]interface{}{data},
		[]string{"category"},
	)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = df.GroupBy(ctx, "category")
	}
}

// BenchmarkGroupByAgg benchmarks grouping with aggregation
func BenchmarkGroupByAgg(b *testing.B) {
	ctx := context.Background()

	// Create dataset with 1000 rows and 10 unique groups
	data := make([]interface{}, 1000)
	for i := range data {
		data[i] = i % 10
	}

	df, _ := New(
		[]interface{}{data},
		[]string{"category"},
	)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		grouped, _ := df.GroupBy(ctx, "category")
		_, _ = grouped.Agg(Count).Show(ctx)
	}
}
