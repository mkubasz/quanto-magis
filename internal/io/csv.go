// Package io provides input/output functionality for reading and writing data files.
package io

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"

	"mkubasz/quanto/internal/dataframe"
)

// Reader provides functionality for reading data from various file formats.
type Reader struct{}

// NewReader creates a new Reader instance.
func NewReader() *Reader {
	return &Reader{}
}

// ReadCSV reads a CSV file and returns a DataFrame containing the data.
// The first row is treated as column headers.
func (r *Reader) ReadCSV(fileName string) (*dataframe.DataFrame, error) {
	file, err := os.Open(fileName) //nolint:gosec // File path is expected to come from user input.
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
func (r *Reader) ReadCSV(fileName string) (_ *dataframe.DataFrame, err error) {
	file, err := os.Open(fileName) //nolint:gosec // File path is expected to come from user input.
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close file: %w", closeErr)
		}
	}()

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV records: %w", err)
	}

	columns, err := createColumns(columnNames, dataRecords)
	if err != nil {
		return nil, fmt.Errorf("failed to create columns: %w", err)
	}

	df, err := dataframe.New(data, columnNames)
	if err != nil {
		return nil, fmt.Errorf("failed to create dataframe: %w", err)
	}

	return df, nil
}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV records: %w", err)
	}

	if len(records) < 2 {
		return nil, errors.New("invalid CSV format: missing data rows")
	}

	columnNames := records[0]
	dataRecords := records[1:]

	columns, err := createColumns(columnNames, dataRecords)
	if err != nil {
		return nil, fmt.Errorf("failed to create columns: %w", err)
	}

	// Convert columns to []interface{} for DataFrame.New
	data := make([]interface{}, len(columns))
	for i, col := range columns {
		data[i] = col.Data
	}

	df, err := dataframe.New(data, columnNames)
	if err != nil {
		return nil, fmt.Errorf("failed to create dataframe: %w", err)
	}

	return df, nil
}

func createColumns(columnNames []string, records [][]string) ([]dataframe.Series[interface{}], error) {
	numColumns := len(columnNames)
	columns := make([]dataframe.Series[interface{}], numColumns)

	for _, record := range records {
		if len(record) != numColumns {
			return nil, errors.New("invalid CSV format: inconsistent number of columns")
		}

		for idx, value := range record {
			f, err := strconv.ParseFloat(value, 64)
			if err != nil {
				columns[idx].Data = append(columns[idx].Data, value)
			} else {
				columns[idx].Data = append(columns[idx].Data, f)
			}
		}
	}

	return columns, nil
}
