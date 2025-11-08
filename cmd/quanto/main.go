// Package main provides the entry point for the Quanto data processing application.
package main

import (
	"log"

	"mkubasz/quanto/internal/cli"
	"mkubasz/quanto/pkg/quanto"
)

func main() {
	// Example usage
	df, err := quanto.NewDataFrame(
		[]interface{}{
			[]interface{}{"A", "B", "A", "D", "E"},
			[]interface{}{1, 2, 3, 4, 5},
		},
		[]string{"col1", "col2"},
	)
	if err != nil {
		log.Printf("Error creating DataFrame: %v\n", err)
		return
	}
	log.Println(df.HasColumn("col1"))

	// Initialize and run CLI
	quantoCli := cli.New()
	if err = quantoCli.Run(); err != nil {
		log.Printf("Error running CLI: %v\n", err)
	}
	log.Println(quantoCli.Session)
}
