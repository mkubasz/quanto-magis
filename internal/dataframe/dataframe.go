// Package dataframe provides a column-oriented data structure for efficient data manipulation.
package dataframe

import (
	"context"
	"fmt"
	"strings"

	"mkubasz/quanto/internal/rdd"
)

// Series represents a single column of homogeneous data.
// Series uses generics for type safety.
type Series[T any] struct {
	Data []T
}

// DataFrame represents a column-oriented data structure,
// similar to a table with named columns.
// Each column is stored as a Series for efficient column-wise operations.
type DataFrame struct {
	series  []Series[interface{}]
	columns []string
	size    int
}

// NewFromRDD creates a new DataFrame from an RDD.
// The DataFrame will have a single column containing all RDD elements.
func NewFromRDD[T any](r *rdd.RDD[T]) *DataFrame {
	var series Series[interface{}]
	for _, v := range r.Collect() {
		series.Data = append(series.Data, v)
	}
	return &DataFrame{
		size:    len(series.Data),
		series:  []Series[interface{}]{series},
		columns: []string{"value"}, // Default column name
	}
}

// New creates a new DataFrame from the provided data and column names.
// Data should be a slice where each element represents a column's data.
//
// Returns ErrInvalidData if data format is invalid.
// Returns ErrInvalidColumnName if column names are empty or don't match data.
func New(data []interface{}, columns []string) (*DataFrame, error) {
	if len(data) == 0 {
		return &DataFrame{
			size:    0,
			series:  []Series[interface{}]{},
			columns: columns,
		}, nil
	}

	if len(columns) == 0 {
		return nil, fmt.Errorf("creating dataframe: %w: must provide column names", ErrInvalidColumnName)
	}

	if len(data) != len(columns) {
		return nil, fmt.Errorf("creating dataframe: %w: data columns (%d) don't match column names (%d)",
			ErrInvalidData, len(data), len(columns))
	}

	// Validate column names
	for i, col := range columns {
		if strings.TrimSpace(col) == "" {
			return nil, fmt.Errorf("creating dataframe: %w: column %d has empty name", ErrInvalidColumnName, i)
		}
	}

	var series []Series[interface{}]
	totalSize := 0

	for i, row := range data {
		rowSlice, ok := row.([]interface{})
		if !ok {
			return nil, fmt.Errorf("creating dataframe: %w: column %d (%s) is not a slice",
				ErrInvalidData, i, columns[i])
		}

		var serie Series[interface{}]
		serie.Data = make([]interface{}, len(rowSlice))
		copy(serie.Data, rowSlice)
		series = append(series, serie)
		totalSize += len(rowSlice)
	}

	return &DataFrame{
		size:    totalSize,
		series:  series,
		columns: columns,
	}, nil
}

// Select returns the Series (column) with the specified name.
//
// Returns ErrColumnNotFound if the column doesn't exist.
// Returns ErrInvalidColumnName if the column name is empty.
func (df *DataFrame) Select(columnName string) (Series[interface{}], error) {
	if strings.TrimSpace(columnName) == "" {
		return Series[interface{}]{}, fmt.Errorf("selecting column: %w: column name is empty", ErrInvalidColumnName)
	}

	idx, err := df.getColumnIndex(columnName)
	if err != nil {
		return Series[interface{}]{}, fmt.Errorf("selecting column %s: %w", columnName, err)
	}

	// Return a copy to ensure immutability
	series := df.series[idx]
	dataCopy := make([]interface{}, len(series.Data))
	copy(dataCopy, series.Data)

	return Series[interface{}]{Data: dataCopy}, nil
}

// getColumnIndex returns the index of a column by name.
// Returns ErrColumnNotFound if the column doesn't exist.
func (df *DataFrame) getColumnIndex(name string) (int, error) {
	for idx, column := range df.columns {
		if column == name {
			return idx, nil
		}
	}
	return -1, ErrColumnNotFound
}

// HasColumn returns true if the DataFrame has a column with the specified name.
func (df *DataFrame) HasColumn(name string) bool {
	for _, column := range df.columns {
		if column == name {
			return true
		}
	}
	return false
}

// Columns returns the names of all columns in the DataFrame.
func (df *DataFrame) Columns() []string {
	cols := make([]string, len(df.columns))
	copy(cols, df.columns)
	return cols
}

// Size returns the total number of elements across all columns.
func (df *DataFrame) Size() int {
	return df.size
}

// NumColumns returns the number of columns in the DataFrame.
func (df *DataFrame) NumColumns() int {
	return len(df.columns)
}

// Distinct returns a new Series containing only unique values.
// The key parameter is ignored and kept for backward compatibility.
//
// Returns ErrEmptyDataFrame if the series is empty.
func (s Series[T]) Distinct(ctx context.Context, key string) (Series[interface{}], error) {
	if err := ctx.Err(); err != nil {
		return Series[interface{}]{}, err
	}

	if len(s.Data) == 0 {
		return Series[interface{}]{}, fmt.Errorf("getting distinct values: %w", ErrEmptyDataFrame)
	}

	uniqueValues := make(map[interface{}]struct{})
	for _, value := range s.Data {
		// Check context periodically for cancellation
		select {
		case <-ctx.Done():
			return Series[interface{}]{}, ctx.Err()
		default:
		}
		uniqueValues[value] = struct{}{}
	}

	distinctValues := make([]interface{}, 0, len(uniqueValues))
	for k := range uniqueValues {
		distinctValues = append(distinctValues, k)
	}

	return Series[interface{}]{Data: distinctValues}, nil
}

// Count returns the number of elements in the Series.
func (s Series[T]) Count() int {
	return len(s.Data)
}

// String returns a string representation of the DataFrame.
func (df *DataFrame) String() string {
	return fmt.Sprintf("DataFrame[rows=%d, columns=%d: %s]",
		df.size, len(df.columns), strings.Join(df.columns, ", "))
}

// String returns a string representation of the Series.
func (s Series[T]) String() string {
	return fmt.Sprintf("Series[count=%d]", len(s.Data))
}
