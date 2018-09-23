package sources

import (
	"context"

	"github.com/shivylp/radium"
)

// NewWikipedia initializes wikipedia based radium source implementation.
func NewWikipedia() *Wikipedia {
	return &Wikipedia{}
}

// Wikipedia implements Source interface using wikipedia for lookups.
type Wikipedia struct {
}

// Search will query en.wikipedia.com to find results and extracts the first
// paragraph of the page.
func (wiki *Wikipedia) Search(ctx context.Context, query radium.Query) ([]radium.Article, error) {
	return nil, nil
}
