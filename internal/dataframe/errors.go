// Package dataframe provides errors used throughout the dataframe package.
package dataframe

import "errors"

// Sentinel errors for common dataframe operations.
var (
	// ErrColumnNotFound is returned when a requested column does not exist.
	ErrColumnNotFound = errors.New("column not found")

	// ErrInvalidColumnName is returned when a column name is empty or invalid.
	ErrInvalidColumnName = errors.New("invalid column name")

	// ErrEmptyDataFrame is returned when an operation requires non-empty data.
	ErrEmptyDataFrame = errors.New("dataframe is empty")

	// ErrInvalidData is returned when input data is malformed.
	ErrInvalidData = errors.New("invalid data format")
)
