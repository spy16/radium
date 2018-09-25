package radium

import (
	"context"
)

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
