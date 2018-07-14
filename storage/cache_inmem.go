package storage

import (
	"sync"

	"github.com/shivylp/radium"
)

// NewDefaultCache initializes an implementartion of radium.Cache
// using in-memeory map
func NewDefaultCache() radium.Cache {
	dc := &defaultCache{}
	dc.mu = new(sync.Mutex)
	dc.data = map[string][]radium.Article{}
	return dc
}

type defaultCache struct {
	mu   *sync.Mutex
	data map[string][]radium.Article
}

func (dc *defaultCache) Search(q radium.Query) ([]radium.Article, error) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	var rs []radium.Article
	if vals, found := dc.data[q.Text]; found {
		rs = append(rs, vals...)
	}
	return rs, nil
}

func (dc *defaultCache) Set(q radium.Query, rs []radium.Article) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	dc.data[q.Text] = rs
	return nil
}
