package radium

import (
	"context"
	"fmt"
)

// New initializes an instance of radium
func New(cache Cache, logger Logger) *Instance {
	ins := &Instance{}
	ins.sources = map[string]Source{}
	ins.cache = cache

	if logger == nil {
		logger = defaultLogger{}
	}
	ins.Logger = logger

	return ins
}

// Instance represents an instance of radium
type Instance struct {
	Logger

	sources map[string]Source
	cache   Cache
}

// RegisterSource adds a new source to the query sources
func (ins Instance) RegisterSource(name string, src Source) error {
	if _, exists := ins.sources[name]; exists {
		return fmt.Errorf("source with given name already exists")
	}

	ins.sources[name] = src
	return nil
}

// GetSources returns a list of registered sources
func (ins Instance) GetSources() map[string]Source {
	return ins.sources
}

// Search using given query and return results if any
func (ins Instance) Search(ctx context.Context, query Query) ([]Article, error) {

	if err := query.Validate(); err != nil {
		return nil, err
	}

	if rs := ins.findInCache(query); rs != nil && len(rs) > 0 {
		return rs, nil
	}

	results := ins.findFromSources(ctx, query)
	go ins.performCaching(query, results)
	return results, nil
}

func (ins Instance) findFromSources(ctx context.Context, query Query) []Article {
	var results []Article
	for srcName, src := range ins.sources {
		resList, err := src.Search(query)
		if err != nil {
			ins.Warnf("source '%s' failed: %s", srcName, err)
			continue
		}

		for _, res := range resList {
			select {
			case <-ctx.Done():
				break
			default:
			}

			if err := res.Validate(); err != nil {
				ins.Warnf("invalid result from source '%s': %s", srcName, err)
				continue
			}

			res.Source = srcName
			results = append(results, res)
		}
	}

	if results == nil {
		results = []Article{}
	}

	return results
}

func (ins Instance) findInCache(query Query) []Article {
	if ins.cache == nil {
		return nil
	}

	rs, err := ins.cache.Search(query)
	if err != nil {
		ins.Warnf("failed to search in cache: %s", err)
		return nil
	}

	return rs
}

func (ins Instance) performCaching(query Query, results []Article) {
	if ins.cache == nil {
		return
	}

	if err := ins.cache.Set(query, results); err != nil {
		ins.Warnf("failed to cache result: %s", err)
	}
}
