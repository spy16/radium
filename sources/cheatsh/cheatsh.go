package cheatsh

import (
	"context"

	"github.com/shivylp/radium"
)

// New initializes a radium.Source implementation using http://cheat.sh as
// the source
func New() *CheatSh {
	csh := &CheatSh{}
	return csh
}

// CheatSh implements radium.Source interface using cheat.sh as the source
type CheatSh struct {
}

// Search performs an HTTP request to http://cheat.sh to find results matching
// the given query.
func (csh CheatSh) Search(ctx context.Context, query radium.Query) ([]radium.Article, error) {
	transformLanguageQuery(&query)

	return executeRequest(ctx, query)
}
