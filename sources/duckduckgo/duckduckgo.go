package duckduckgo

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/spby/radium"
)

const ddgURL = "https://api.duckduckgo.com"

// New initializes a duckduckgo based radium source implementation.
func New() *DuckDuckGo {
	return &DuckDuckGo{}
}

// DuckDuckGo implements radium.Source interface using DuckDuckGo search
// engine project.
type DuckDuckGo struct {
}

// Search makes request to duckduckgo instant answers API to fetch the abstract.
func (ddg DuckDuckGo) Search(ctx context.Context, query radium.Query) ([]radium.Article, error) {
	req, err := makeDDGRequest(ctx, query)
	if err != nil {
		return nil, err
	}

	cl := http.Client{}

	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ires := instantResult{}
	if err := json.NewDecoder(resp.Body).Decode(&ires); err != nil {
		return nil, err
	}

	result := radium.Article{}
	result.Title = query.Text
	result.Content = ires.Abstract
	return []radium.Article{result}, nil
}

func makeDDGRequest(ctx context.Context, query radium.Query) (*http.Request, error) {
	pu, err := url.Parse(ddgURL)
	if err != nil {
		return nil, err
	}

	params := pu.Query()
	params.Set("q", query.Text)
	params.Set("format", "json")
	pu.RawQuery = params.Encode()

	req, err := http.NewRequest(http.MethodGet, pu.String(), nil)
	if err != nil {
		return nil, err
	}
	req.WithContext(ctx)

	return req, nil
}

type instantResult struct {
	Abstract    string `json:"AbstractText"`
	AbstractURL string `json:"AbstractURL"`
}
