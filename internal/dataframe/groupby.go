package dataframe

import (
	"context"
	"fmt"
)

// DataFrameGroupBy represents a grouped DataFrame for aggregation operations.
type DataFrameGroupBy struct {
	df         *DataFrame
	columnName string
	aggs       []func([]interface{}) int
	groups     map[interface{}][]interface{}
}

// GroupBy creates a grouped DataFrame based on the specified column.
// Subsequent aggregation operations can be performed on the groups.
//
// Returns ErrColumnNotFound if the column doesn't exist.
// Returns ErrInvalidColumnName if the column name is empty.
func (df *DataFrame) GroupBy(ctx context.Context, name string) (*DataFrameGroupBy, error) {
	// Check context
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if !df.HasColumn(name) {
		return nil, fmt.Errorf("grouping by column %s: %w", name, ErrColumnNotFound)
	}

	// Get column index
	index, err := df.getColumnIndex(name)
	if err != nil {
		return nil, fmt.Errorf("grouping by column %s: %w", name, err)
	}

	// Create groups based on unique values in the column
	groups := make(map[interface{}][]interface{})
	for _, el := range df.series[index].Data {
		// Check context periodically
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		groups[el] = append(groups[el], el)
	}

	return &DataFrameGroupBy{
		df:         df,
		columnName: name,
		groups:     groups,
	}, nil
}

// Agg adds an aggregation function to be applied to each group.
// The function takes a slice of values and returns an integer result.
// Multiple aggregation functions can be chained.
func (dfg *DataFrameGroupBy) Agg(f func([]interface{}) int) *DataFrameGroupBy {
	dfg.aggs = append(dfg.aggs, f)
	return dfg
}

// Show materializes the grouped DataFrame with aggregations applied.
// Returns a new DataFrame with columns for the grouping key and aggregation results.
//
// Returns ErrInvalidData if Show is called before any aggregations are added.
func (dfg *DataFrameGroupBy) Show(ctx context.Context) (*DataFrame, error) {
	// Check context
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if len(dfg.aggs) == 0 {
		return nil, fmt.Errorf("showing grouped data: %w: no aggregation functions specified", ErrInvalidData)
	}

	hash := make([]interface{}, 0)
	groups := make([]interface{}, 0)

	for key, group := range dfg.groups {
		// Check context periodically
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		for _, agg := range dfg.aggs {
			hash = append(hash, key)
			groups = append(groups, agg(group))
		}
	}

	return New([]interface{}{hash, groups}, []string{dfg.columnName, "count"})
}

// Count is an aggregation function that returns the number of elements in a group.
func Count(group []interface{}) int {
	return len(group)
}

// Sum is an aggregation function that returns the sum of numeric elements in a group.
// Non-numeric values are ignored.
func Sum(group []interface{}) int {
	sum := 0
	for _, val := range group {
		if num, ok := val.(int); ok {
			sum += num
		}
	}
	return sum
}
