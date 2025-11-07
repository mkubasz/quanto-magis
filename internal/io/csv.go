package io

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"mkubasz/quanto/internal/dataframe"
)

type Reader struct{}

func NewReader() *Reader {
	return &Reader{}
}

func (r *Reader) ReadCSV(fileName string) (*dataframe.DataFrame, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV records: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("invalid CSV format: missing data rows")
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

	return dataframe.New(data, columnNames), nil
}

func createColumns(columnNames []string, records [][]string) ([]dataframe.Series[interface{}], error) {
	numColumns := len(columnNames)
	columns := make([]dataframe.Series[interface{}], numColumns)

	for _, record := range records {
		if len(record) != numColumns {
			return nil, fmt.Errorf("invalid CSV format: inconsistent number of columns")
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
