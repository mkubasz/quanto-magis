package main

type DataFrameGroupBy struct {
	df *DataFrame
}

func (df *DataFrame) GroupBy(name string) *DataFrameGroupBy {
	return &DataFrameGroupBy{df: df}
}
