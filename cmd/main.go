package main

import (
	"fmt"
	"mkubasz/quanto/quantoCore"
)

func main() {
	c := quantoCore.NewDataFrame([]interface{}{[]interface{}{"A", "B", "A", "D", "E"}, []interface{}{1, 2, 3, 4, 5}}, []string{"col1", "col2"})
	fmt.Println(c.HasColumn("col1"))
}
