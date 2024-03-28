package quantoCore

import "errors"

type DataFrameGroupBy struct {
	df         *DataFrame
	columnName string
	aggs       []func([]interface{}) int
	groups     map[interface{}][]interface{}
}

func (df *DataFrame) GroupBy(name string) (*DataFrameGroupBy, error) {
	if !df.HasColumn(name) {
		return nil, errors.New("column not found")
	}
	// Create a group by object
	index, _ := df.getColumnIndex(name)
	groups := make(map[interface{}][]interface{})
	for _, el := range df.series[index].data {
		groups[el] = append(groups[el], el)
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

func (dfg *DataFrameGroupBy) Agg(f func([]interface{}) int) *DataFrameGroupBy {
	dfg.aggs = append(dfg.aggs, f)
	return dfg
}

func (dfg *DataFrameGroupBy) Show() *DataFrame {
	hash := make([]interface{}, 0)
	groups := make([]interface{}, 0)
	for key, group := range dfg.groups {
		for _, agg := range dfg.aggs {
			hash = append(hash, key)
			groups = append(groups, agg(group))
		}
	}
	return NewDataFrame([]interface{}{
		hash,
		groups,
	}, []string{dfg.columnName, "count"})
}

func Count(group []interface{}) int {
	return len(group)
}
