package radium

import (
	"context"
	"sync"
)

// Strategy implementation is responsible for performing queries
// against given set of sources using a particular approach.
type Strategy interface {
	Execute(ctx context.Context, query Query, sources []RegisteredSource) ([]Article, error)
}

// NewConcurrent initializes a concurrent radium strategy
func NewConcurrent(logger Logger) *Concurrent {
	return &Concurrent{
		Logger: logger,
	}
}

// Concurrent is a radium strategy implementation.
type Concurrent struct {
	Logger
}

// Execute the query against given list of sources concurrently. This strategy
// ingores the source errors and simply logs them.
func (con Concurrent) Execute(ctx context.Context, query Query, sources []RegisteredSource) ([]Article, error) {
	results := newSafeResults()
	wg := &sync.WaitGroup{}

	for _, source := range sources {
		wg.Add(1)

		go func(wg *sync.WaitGroup, src RegisteredSource, rs *safeResults) {
			srcResults, err := src.Search(ctx, query)
			if err != nil {
				con.Warnf("source '%s' failed: %s", src.Name, err)
				return
			}

			rs.extend(srcResults, src.Name, con.Logger)
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
	results []Article
}

func (sr *safeResults) extend(results []Article, srcName string, logger Logger) {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	for _, res := range results {
		if err := res.Validate(); err != nil {
			logger.Warnf("ignoring invalid result from source '%s': %s", srcName, err)
			continue
		}
		res.Source = srcName

		sr.results = append(sr.results, res)
	}

}

// NewNthResult initializes NthResult strategy with given n
func NewNthResult(n int, logger Logger) *NthResult {
	return &NthResult{stopAt: n, Logger: logger}
}

// NthResult implements a radium search strategy. This strategy
// executes search in the given order of sources and stops at nth
// result or if all the sources are executed.
type NthResult struct {
	Logger

	stopAt int
}

// Execute each source in srcs until n results are obtained or all sources have
// been executed. This strategy returns on first error.
func (nth *NthResult) Execute(ctx context.Context, query Query, srcs []RegisteredSource) ([]Article, error) {
	results := []Article{}
	for _, src := range srcs {
		select {
		case <-ctx.Done():
			break
		default:
		}

		srcResults, err := src.Search(ctx, query)
		if err != nil {
			return nil, err
		}

		for _, res := range srcResults {
			if err := res.Validate(); err != nil {
				nth.Warnf("ignoring invalid result  from '%s': %s", src.Name, err)
				continue
			}

			res.Source = src.Name
			results = append(results, res)
		}

		if len(results) >= nth.stopAt {
			break
		}
	}
	return results, nil
}
