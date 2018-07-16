package radium

import (
	"fmt"
	"strings"
)

// Article represents a radium article which can be a
// snippet, wiki, answer etc.
type Article struct {
	// Title should represent a short title or a summary of the
	// content
	Title string `json:"title"`

	// Content should contain the body of the article which may
	// be a snippet, wiki article, answer etc.
	Content string `json:"content"`

	// ContentType should contain the type of the content returned
	// so that radium can use that info to render the content. It
	// should be one of markdown, json, yaml, html
	ContentType string `json:"content_type"`

	// Attribs can contain type of the article, keywords etc.
	Attribs map[string]string `json:"attribs"`

	// License can contain name of the license if applicable
	License string `json:"license"`

	// Source should contain the name of the registered source
	// which returned this article. This will be added automatically
	// by radium.
	Source string `json:"source"`
}

// Validate the article model
func (article Article) Validate() error {
	if strings.TrimSpace(article.Title) == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if strings.TrimSpace(article.Content) == "" {
		return fmt.Errorf("content cannot be empty")
	}
	return nil
}

// Query represents a user query for a post
type Query struct {
	// Text is the primary search criteria. Sources must
	// use this to find relevant results
	Text string `json:"text"`

	// Attribs can be used by sources to further filter down
	// the results
	Attribs map[string]string `json:"attribs"`
}

// Validate for empty or invalid queries
func (query Query) Validate() error {
	if strings.TrimSpace(query.Text) == "" {
		return fmt.Errorf("invalid query")
	}
	return nil
}
