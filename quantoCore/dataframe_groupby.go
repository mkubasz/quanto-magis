package main

import "errors"

type DataFrameGroupBy struct {
	df         *DataFrame
	columnName string
	aggs       []func() *DataFrame
	groups     map[interface{}][]int
}

func (df *DataFrame) GroupBy(name string) (*DataFrameGroupBy, error) {
	if !df.HasColumn(name) {
		return nil, errors.New("column not found")
	}

	// Create a group by object
	index, _ := df.getColumnIndex(name)
	groups := make(map[interface{}][]int)
	for idx, el := range df.series[index].data {
		groups[el] = append(groups[el], idx)
	}
	return &DataFrameGroupBy{
		df:         df,
		columnName: name,
		groups:     groups,
	}, nil
}

func (df *DataFrame) HasColumn(name string) bool {
	for _, column := range df.columns {
		if column == name {
			return true
		}
	}
	return false
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
