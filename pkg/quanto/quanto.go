package quanto

import (
	"mkubasz/quanto/internal/dataframe"
	"mkubasz/quanto/internal/io"
	"mkubasz/quanto/internal/rdd"
	"mkubasz/quanto/internal/session"
)

// Session wraps the internal QuantoSession with additional functionality
type Session struct {
	*session.QuantoSession
	Reader *io.Reader
}

// NewSession creates a new Quanto session with all capabilities
func NewSession() *Session {
	return &Session{
		QuantoSession: session.New(),
		Reader:        io.NewReader(),
	}
}

// Parallelize creates an RDD from a slice of data
func (s *Session) Parallelize(data []interface{}) *rdd.RDD[interface{}] {
	return rdd.New(data)
}

// ReadCSV reads a CSV file and returns a DataFrame
func (s *Session) ReadCSV(fileName string) (*dataframe.DataFrame, error) {
	return s.Reader.ReadCSV(fileName)
}

// Re-export commonly used types for convenience
type (
	DataFrame        = dataframe.DataFrame
	DataFrameGroupBy = dataframe.DataFrameGroupBy
	Mode             = session.Mode
)

// Re-export commonly used functions
var (
	NewDataFrame = dataframe.New
	Count        = dataframe.Count
)

// NewDataFrameFromRDD creates a DataFrame from an RDD
func NewDataFrameFromRDD[T any](r *rdd.RDD[T]) *dataframe.DataFrame {
	return dataframe.NewFromRDD(r)
}

// NewRDD creates a new RDD from a slice of data
func NewRDD[T any](data []T) *rdd.RDD[T] {
	return rdd.New(data)
}

// Re-export mode constants
const (
	Local   = session.Local
	Cluster = session.Cluster
)
