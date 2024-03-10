package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

type Read struct{}

func (r *Read) csv(fileName string) {
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

	for _, record := range records {
		fmt.Println(record)
	}
}
