package radium

import (
	"context"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

// Default registered strategies
const (
	Strategy1st        = "1st"
	StrategyConcurrent = "concurrent"
)

// New initializes an instance of radium
func New(cache Cache, logger Logger) *Instance {
	if cache == nil {
		// only useful in case of server or clipboard mode
		cache = &defaultCache{
			mu:   &sync.Mutex{},
			data: map[string][]Article{},
		}
	}

	if logger == nil {
		logger = defaultLogger{
			logger: logrus.New(),
		}
	}

	ins := &Instance{}
	ins.cache = cache
	ins.Logger = logger
	ins.strategies = map[string]Strategy{
		Strategy1st:        NewNthResult(1, ins.Logger),
		StrategyConcurrent: NewConcurrent(ins.Logger),
	}

	return ins
}

// Instance represents an instance of radium
type Instance struct {
	Logger

	sources    []RegisteredSource
	strategies map[string]Strategy
	cache      Cache
}

// RegisterStrategy adds a new source to the query sources
func (ins *Instance) RegisterStrategy(name string, strategy Strategy) {
	if ins.strategies == nil {
		ins.strategies = map[string]Strategy{}
	}

	ins.strategies[name] = strategy
}

// RegisterSource adds a new source to the query sources
func (ins *Instance) RegisterSource(name string, src Source) error {
	for _, entry := range ins.sources {
		if name == entry.Name {
			return fmt.Errorf("source with given name already exists")
		}
	}

	ins.sources = append(ins.sources, RegisteredSource{
		Name:   name,
		Source: src,
	})
	return nil
}

// GetSources returns a list of registered sources
func (ins Instance) GetSources() []RegisteredSource {
	return ins.sources
}

// Search using given query and return results if any
func (ins Instance) Search(ctx context.Context, query Query, strategyName string) ([]Article, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}

	if rs := ins.findInCache(query); rs != nil && len(rs) > 0 {
		ins.Infof("cache hit for '%s'", query.Text)
		return rs, nil
	}
	ins.Infof("cache miss for '%s'", query.Text)

	strategy, exists := ins.strategies[strategyName]
	if !exists {
		return nil, fmt.Errorf("no such strategy: %s", strategyName)
	}

	results, err := strategy.Execute(ctx, query, ins.sources)
	if err != nil {
		return nil, err
	}

	go ins.performCaching(query, results)
	return results, nil
}

func (ins Instance) findInCache(query Query) []Article {
	if ins.cache == nil {
		return nil
	}

	rs, err := ins.cache.Search(context.Background(), query)
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

// Source implementation is responsible for providing
// external data source to query for results.
type Source interface {
	Search(ctx context.Context, q Query) ([]Article, error)
}

// RegisteredSource embeds given Source along with the registered name.
type RegisteredSource struct {
	Name string
	Source
}
