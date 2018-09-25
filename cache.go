package radium

import (
	"context"
	"sync"
)

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

type defaultCache struct {
	mu   *sync.Mutex
	data map[string][]Article
}

func (dc *defaultCache) Search(ctx context.Context, q Query) ([]Article, error) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	var rs []Article
	if vals, found := dc.data[q.Text]; found {
		rs = append(rs, vals...)
	}
	return rs, nil
}

func (dc *defaultCache) Set(q Query, rs []Article) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	dc.data[q.Text] = rs
	return nil
}
