// Package quanto provides a unified API for data processing with DataFrames and RDDs.
package quanto

import (
	"mkubasz/quanto/internal/dataframe"
	"mkubasz/quanto/internal/io"
	"mkubasz/quanto/internal/rdd"
	"mkubasz/quanto/internal/session"
)

// Session wraps the internal QuantoSession with additional functionality.
type Session struct {
	*session.QuantoSession
	Reader *io.Reader
}

// NewSession creates a new Quanto session with all capabilities.
func NewSession() *Session {
	return &Session{
		QuantoSession: session.New(),
		Reader:        io.NewReader(),
	}
}

// Parallelize creates an RDD from a slice of data.
func (s *Session) Parallelize(data []interface{}) *rdd.RDD[interface{}] {
	return rdd.New(data)
}

// ReadCSV reads a CSV file and returns a DataFrame.
func (s *Session) ReadCSV(fileName string) (*dataframe.DataFrame, error) {
	return s.Reader.ReadCSV(fileName)
}

// Re-export commonly used types for convenience

// DataFrame is an alias for dataframe.DataFrame providing column-oriented data structure.
type DataFrame = dataframe.DataFrame

// GroupBy is an alias for dataframe.GroupBy for grouped aggregation operations.
type GroupBy = dataframe.GroupBy

// Mode is an alias for session.Mode representing execution modes.
type Mode = session.Mode

// NewDataFrame creates a new DataFrame from columns and column names.
func NewDataFrame(columns []interface{}, columnNames []string) (*dataframe.DataFrame, error) {
	return dataframe.New(columns, columnNames)
}

// Count is a dataframe aggregation function that counts elements.
func Count(values []interface{}) int {
	return dataframe.Count(values)
}

// NewDataFrameFromRDD creates a DataFrame from an RDD.
func NewDataFrameFromRDD[T any](r *rdd.RDD[T]) *dataframe.DataFrame {
	return dataframe.NewFromRDD(r)
}

// NewRDD creates a new RDD from a slice of data.
func NewRDD[T any](data []T) *rdd.RDD[T] {
	return rdd.New(data)
}

// Re-export mode constants.
const (
	Local   = session.Local
	Cluster = session.Cluster
)
