package main

type DataFrameGroupBy struct {
	df         *DataFrame
	columnName string
	aggs       []func() *DataFrame
}

func (df *DataFrame) GroupBy(name string) *DataFrameGroupBy {
	return &DataFrameGroupBy{
		df:         df,
		columnName: name,
	}
}

func (dfg *DataFrameGroupBy) Agg(f func() *DataFrame) *DataFrameGroupBy {
	dfg.aggs = append(dfg.aggs, f)
	return dfg
}

func (dfg *DataFrameGroupBy) Show() *DataFrame {
	return nil
}

func (df *DataFrame) Count() *DataFrame {
	return nil
}
