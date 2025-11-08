package dataframe_test

import (
	"context"
	"errors"
	"testing"

	"mkubasz/quanto/internal/dataframe"
	"mkubasz/quanto/internal/rdd"
)

// TestNewFromRDD verifies DataFrame creation from RDD.
func TestNewFromRDD(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		wantSize int
		wantCols int
	}{
		{
			name:     "create from string array",
			input:    []interface{}{"A", "B", "C", "D", "E"},
			wantSize: 5,
			wantCols: 1,
		},
		{
			name:     "create from integer array",
			input:    []interface{}{1, 2, 3},
			wantSize: 3,
			wantCols: 1,
		},
		{
			name:     "create from empty array",
			input:    []interface{}{},
			wantSize: 0,
			wantCols: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := rdd.New(tt.input)
			df := dataframe.NewFromRDD(r)

			if df == nil {
				t.Error("expected non-nil DataFrame")
				return
			}

			if df.Size() != tt.wantSize {
				t.Errorf("size = %d, want %d", df.Size(), tt.wantSize)
			}

			if df.NumColumns() != tt.wantCols {
				t.Errorf("columns = %d, want %d", df.NumColumns(), tt.wantCols)
			}
		})
	}
}

// TestNew verifies DataFrame creation with validation.
//
//nolint:funlen // Test functions require comprehensive test cases.
func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		data     []interface{}
		columns  []string
		wantErr  error
		wantSize int
	}{
		{
			name: "create with two columns",
			data: []interface{}{
				[]interface{}{"A", "B", "C"},
				[]interface{}{1, 2, 3},
			},
			columns:  []string{"col1", "col2"},
			wantErr:  nil,
			wantSize: 6,
		},
		{
			name:     "empty data with columns",
			data:     []interface{}{},
			columns:  []string{"col1"},
			wantErr:  nil,
			wantSize: 0,
		},
		{
			name: "missing column names",
			data: []interface{}{
				[]interface{}{1, 2, 3},
			},
			columns: []string{},
			wantErr: dataframe.ErrInvalidColumnName,
		},
		{
			name: "column count mismatch",
			data: []interface{}{
				[]interface{}{1, 2, 3},
				[]interface{}{4, 5, 6},
			},
			columns: []string{"col1"},
			wantErr: dataframe.ErrInvalidData,
		},
		{
			name: "empty column name",
			data: []interface{}{
				[]interface{}{1, 2, 3},
			},
			columns: []string{""},
			wantErr: dataframe.ErrInvalidColumnName,
		},
		{
			name: "whitespace column name",
			data: []interface{}{
				[]interface{}{1, 2, 3},
			},
			columns: []string{"  "},
			wantErr: dataframe.ErrInvalidColumnName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df, err := dataframe.New(tt.data, tt.columns)

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

			if df.Size() != tt.wantSize {
				t.Errorf("size = %d, want %d", df.Size(), tt.wantSize)
			}
		})
	}
}

// TestSelect verifies column selection.
//
//nolint:funlen // Test functions require comprehensive test cases.
func TestSelect(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() *dataframe.DataFrame
		columnName string
		wantErr    error
		wantCount  int
	}{
		{
			name: "select existing column",
			setup: func() *dataframe.DataFrame {
				df, _ := dataframe.New(
					[]interface{}{
						[]interface{}{"A", "B", "C"},
						[]interface{}{1, 2, 3},
					},
					[]string{"col1", "col2"},
				)
				return df
			},
			columnName: "col1",
			wantErr:    nil,
			wantCount:  3,
		},
		{
			name: "select non-existing column",
			setup: func() *dataframe.DataFrame {
				df, _ := dataframe.New(
					[]interface{}{
						[]interface{}{"A", "B", "C"},
					},
					[]string{"col1"},
				)
				return df
			},
			columnName: "col2",
			wantErr:    dataframe.ErrColumnNotFound,
		},
		{
			name: "select with empty name",
			setup: func() *dataframe.DataFrame {
				df, _ := dataframe.New(
					[]interface{}{
						[]interface{}{1, 2, 3},
					},
					[]string{"col1"},
				)
				return df
			},
			columnName: "",
			wantErr:    dataframe.ErrInvalidColumnName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df := tt.setup()
			series, err := df.Select(tt.columnName)

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

			if series.Count() != tt.wantCount {
				t.Errorf("count = %d, want %d", series.Count(), tt.wantCount)
			}
		})
	}
}

// TestHasColumn verifies column existence checking.
func TestHasColumn(t *testing.T) {
	df, _ := dataframe.New(
		[]interface{}{
			[]interface{}{"A", "B", "C"},
			[]interface{}{1, 2, 3},
		},
		[]string{"col1", "col2"},
	)

	tests := []struct {
		name       string
		columnName string
		want       bool
	}{
		{"existing column", "col1", true},
		{"another existing column", "col2", true},
		{"non-existing column", "col3", false},
		{"empty name", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := df.HasColumn(tt.columnName)
			if got != tt.want {
				t.Errorf("HasColumn(%s) = %v, want %v", tt.columnName, got, tt.want)
			}
		})
	}
}

// TestDistinct verifies distinct value extraction.
func TestDistinct(t *testing.T) {
	tests := []struct {
		name      string
		input     []interface{}
		wantCount int
		wantErr   error
	}{
		{
			name:      "distinct from duplicates",
			input:     []interface{}{"A", "B", "A", "C", "B"},
			wantCount: 3, // A, B, C
			wantErr:   nil,
		},
		{
			name:      "all unique values",
			input:     []interface{}{"A", "B", "C"},
			wantCount: 3,
			wantErr:   nil,
		},
		{
			name:      "all same values",
			input:     []interface{}{"A", "A", "A"},
			wantCount: 1,
			wantErr:   nil,
		},
		{
			name:      "empty series",
			input:     []interface{}{},
			wantCount: 0,
			wantErr:   dataframe.ErrEmptyDataFrame,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			series := dataframe.Series[interface{}]{Data: tt.input}

			distinct, err := series.Distinct(ctx, "")

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

			if distinct.Count() != tt.wantCount {
				t.Errorf("count = %d, want %d", distinct.Count(), tt.wantCount)
			}
		})
	}
}

// TestDistinctCancellation verifies context cancellation.
func TestDistinctCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately.

	// Create large dataset.
	data := make([]interface{}, 10000)
	for i := range data {
		data[i] = i % 100 // Create some duplicates.
	}

	series := dataframe.Series[interface{}]{Data: data}
	_, err := series.Distinct(ctx, "")

	if err == nil {
		t.Error("expected error due to context cancellation")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled error, got: %v", err)
	}
}

// TestColumns verifies column name retrieval.
func TestColumns(t *testing.T) {
	df, _ := dataframe.New(
		[]interface{}{
			[]interface{}{1, 2, 3},
			[]interface{}{4, 5, 6},
		},
		[]string{"col1", "col2"},
	)

	cols := df.Columns()

	if len(cols) != 2 {
		t.Errorf("columns length = %d, want 2", len(cols))
	}

	expected := map[string]bool{"col1": true, "col2": true}
	for _, col := range cols {
		if !expected[col] {
			t.Errorf("unexpected column: %s", col)
		}
	}
}

// TestDataFrameImmutability verifies that operations don't modify original data.
func TestDataFrameImmutability(t *testing.T) {
	original := []interface{}{1, 2, 3}
	df, _ := dataframe.New(
		[]interface{}{original},
		[]string{"col1"},
	)

	// Select should return a copy.
	series, _ := df.Select("col1")

	// Modify the returned series.
	series.Data[0] = 999

	// Verify original DataFrame is unchanged.
	series2, _ := df.Select("col1")
	if series2.Data[0] == 999 {
		t.Error("DataFrame was mutated, expected immutability")
	}
}

// BenchmarkSelect benchmarks column selection.
func BenchmarkSelect(b *testing.B) {
	data := make([]interface{}, 10)
	for i := 0; i < 10; i++ {
		col := make([]interface{}, 1000)
		for j := range col {
			col[j] = j
		}
		data[i] = col
	}

	columns := make([]string, 10)
	for i := range columns {
		columns[i] = string(rune('a' + i))
	}

	df, _ := dataframe.New(data, columns)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = df.Select("e")
	}
}

// BenchmarkDistinct benchmarks distinct value extraction.
func BenchmarkDistinct(b *testing.B) {
	ctx := context.Background()

	data := make([]interface{}, 10000)
	for i := range data {
		data[i] = i % 100 // 100 unique values
	}

	series := dataframe.Series[interface{}]{Data: data}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = series.Distinct(ctx, "")
	}
}
