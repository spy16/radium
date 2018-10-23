package wikipedia

import (
	"context"

	"github.com/spy16/radium"
)

// New initializes wikipedia based radium source implementation.
func New(mwURL string) *Wikipedia {
	return &Wikipedia{
		url: mwURL,
	}
}

// Wikipedia implements Source interface using wikipedia for lookups.
type Wikipedia struct {
	url string
}

// Search will query en.wikipedia.com to find results and extracts the first
// paragraph of the page.
func (wiki *Wikipedia) Search(ctx context.Context, query radium.Query) ([]radium.Article, error) {
	lang, found := query.Attribs["language"]
	if !found {
		lang = "en"
	}

	req, err := NewRequest(wiki.url, query.Text, lang)
	if err != nil {
		return nil, err
	}

	resp, err := req.Execute(ctx, true)
	if err != nil {
		return nil, err
	}

	page, err := resp.Page()
	if err != nil {
		return nil, err
	}

	return pageToRadiumArticle(page), nil
}

func pageToRadiumArticle(page *Page) []radium.Article {
	article := radium.Article{}
	article.Title = page.Title
	article.Content = page.Content
	article.License = "CC BY-SA 3.0 Unported License"

	return []radium.Article{article}
}
