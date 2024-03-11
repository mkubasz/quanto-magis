package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

type Read struct{}

func (r *Read) Csv(fileName string) *DataFrame {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	columnNames := records[0]
	size := len(columnNames)
	columns := make([]Series[interface{}], size)

	for _, record := range records[1:] {
		for idx, value := range record {
			f, err := strconv.ParseFloat(value, 64)
			if err != nil {
				columns[idx].data = append(columns[idx].data, value)
			} else {
				columns[idx].data = append(columns[idx].data, f)
			}
		}
	}

	return &DataFrame{
		series: columns,
		columns:  columnNames,
		size:     len(records) - 1,
	}
}
