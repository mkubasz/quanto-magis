package main

import (
	"fmt"

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
		fmt.Printf("Error creating DataFrame: %v\n", err)
		return
	}
	fmt.Println(df.HasColumn("col1"))

	// Initialize and run CLI
	quantoCli := cli.New()
	quantoCli.Run()
	fmt.Println(quantoCli.Session)
}
