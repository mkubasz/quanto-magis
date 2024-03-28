package quantoCore

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type Read struct{}

func (r *Read) Csv(fileName string) (*DataFrame, error) {
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

	return &DataFrame{
		series:  columns,
		columns: columnNames,
		size:    len(dataRecords),
	}, nil
}

func createColumns(columnNames []string, records [][]string) ([]Series[interface{}], error) {
	numColumns := len(columnNames)
	columns := make([]Series[interface{}], numColumns)

	for _, record := range records {
		if len(record) != numColumns {
			return nil, fmt.Errorf("invalid CSV format: inconsistent number of columns")
		}

		for idx, value := range record {
			f, err := strconv.ParseFloat(value, 64)
			if err != nil {
				columns[idx].data = append(columns[idx].data, value)
			} else {
				columns[idx].data = append(columns[idx].data, f)
			}
		}
	}

	return columns, nil
}
