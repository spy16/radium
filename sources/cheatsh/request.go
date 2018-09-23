package cheatsh

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/shivylp/radium"
)

const cheatShURL = "http://cheat.sh"

func executeRequest(ctx context.Context, query radium.Query) ([]radium.Article, error) {
	req, err := makeRequest(ctx, query)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return makeResponse(query, data), nil
}

func makeResponse(query radium.Query, data []byte) []radium.Article {
	result := radium.Article{}
	result.Content = string(data)
	result.ContentType = "plaintext"
	result.Title = query.Text
	result.Attribs = query.Attribs

	return []radium.Article{result}
}

func makeRequest(ctx context.Context, query radium.Query) (*http.Request, error) {
	pu, err := url.Parse(cheatShURL)
	if err != nil {
		return nil, err
	}
	queryParams := pu.Query()
	queryStr := url.PathEscape(strings.Replace(strings.TrimSpace(query.Text), " ", "+", -1))

	if lang, found := query.Attribs["language"]; found {
		appendPath(pu, url.PathEscape(lang), queryStr)
	} else {
		appendPath(pu, queryStr)
	}

	if _, found := query.Attribs["color"]; !found {
		queryParams.Set("T", "")
	}

	pu.RawQuery = queryParams.Encode()
	req, err := http.NewRequest(http.MethodGet, pu.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "curl/7.54.0")
	req.WithContext(ctx)

	return req, nil
}

func isTrue(s string) bool {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "yes", "y", "yea", "t", "true":
		return true
	default:
		return false
	}
}

func appendPath(pu *url.URL, segs ...string) {
	segs = append(segs, pu.Path)
	p := filepath.Join(segs...)
	pu.Path = p
}
