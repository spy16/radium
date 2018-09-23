package strategies

import (
	"context"

	"github.com/shivylp/radium"
)

// NewNthResult initializes NthResult strategy with given n
func NewNthResult(n int, logger radium.Logger) *NthResult {
	return &NthResult{stopAt: n, Logger: logger}
}

// NthResult implements a radium search strategy. This strategy
// executes search in the given order of sources and stops at nth
// result or if all the sources are executed.
type NthResult struct {
	radium.Logger

	stopAt int
}

// Execute each source in srcs until n results are obtained or all sources have
// been executed. This strategy returns on first error.
func (nth *NthResult) Execute(ctx context.Context, query radium.Query, srcs []radium.RegisteredSource) ([]radium.Article, error) {
	results := []radium.Article{}
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

		results = append(results, srcResults...)
		if len(results) >= nth.stopAt {
			break
		}
	}
	return results, nil
}
