package strategies

import (
	"context"
	"sync"

	"github.com/shivylp/radium"
)

// NewConcurrent initializes a concurrent radium strategy
func NewConcurrent(logger radium.Logger) *Concurrent {
	return &Concurrent{
		Logger: logger,
	}
}

// Concurrent is a radium strategy implementation.
type Concurrent struct {
	radium.Logger
}

// Execute the query against given list of sources concurrently. This strategy
// ingores the source errors and simply logs them.
func (con Concurrent) Execute(ctx context.Context, query radium.Query, sources []radium.RegisteredSource) ([]radium.Article, error) {
	results := newSafeResults()
	wg := &sync.WaitGroup{}

	for _, source := range sources {
		wg.Add(1)

		go func(wg *sync.WaitGroup, src radium.RegisteredSource, rs *safeResults) {
			srcResults, err := src.Search(ctx, query)
			if err != nil {
				con.Warnf("source '%s' failed: %s", src.Name, err)
				return
			}

			rs.extend(src.Name, con.Logger, srcResults)
			wg.Done()
		}(wg, source, results)
	}

	wg.Wait()
	return results.results, nil
}

func newSafeResults() *safeResults {
	return &safeResults{
		mu: &sync.Mutex{},
	}
}

type safeResults struct {
	mu      *sync.Mutex
	results []radium.Article
}

func (sr *safeResults) extend(results []radium.Article, srcName string, logger radium.Logger) {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	for _, res := range results {
		if err := res.Validate(); err != nil {
			logger.Warnf("ignoring invalid result from source '%s': %s", srcName, err)
			continue
		}

		sr.results = append(sr.results, res)
	}

}
