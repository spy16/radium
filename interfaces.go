package radium

import (
	"fmt"
	"log"
)

// Source implementation is responsible for providing
// external data source to query for results.
type Source interface {
	Search(q Query) ([]Article, error)
}

// Logger implementation should provide logging
// functionality to the radium instance. Log levels
// should be managed externally.
type Logger interface {
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// Cache implementation is responsible for caching
// a given query-results pair for later use
type Cache interface {
	Source

	// Set should store the given pair in a caching
	// backend for fast access. If an entry with same
	// query already exists, it should be replaced
	// with the new results slice
	Set(q Query, rs []Article) error
}

// defaultLogger implements Logger using log package
type defaultLogger struct {
}

func (dl defaultLogger) Infof(format string, args ...interface{}) {
	log.Printf("INFO : %s", fmt.Sprintf(format, args...))
}

func (dl defaultLogger) Warnf(format string, args ...interface{}) {
	log.Printf("WARN : %s", fmt.Sprintf(format, args...))
}

func (dl defaultLogger) Errorf(format string, args ...interface{}) {
	log.Printf("ERR  : %s", fmt.Sprintf(format, args...))
}
