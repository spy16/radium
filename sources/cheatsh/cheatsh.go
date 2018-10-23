package cheatsh

import (
	"context"
	"strings"

	"github.com/spby/radium"
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

	raw, err := executeRequest(ctx, query)
	if err != nil {
		return nil, err
	}

	results := []radium.Article{}
	for _, res := range raw {
		if !strings.HasPrefix(res.Content, "Unknown topic") {
			results = append(results, res)
		}
	}

	return results, nil
}
