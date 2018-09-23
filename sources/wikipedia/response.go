/*
Copyright (c) 2015 Fredrik Wallgren

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE
*/

package wikipedia

import (
	"net/url"
	"strings"
	"time"
)

// Page contains the parsed data.
type Page struct {
	ID       int
	Title    string
	Content  string
	Language string
	URL      string
	Redirect *redirect
}

// Response contains the raw data the API returns.
type Response struct {
	Batchcomplete string
	Query         query
}

// Page parses the raw data and returns a Page with the relevant data.
func (r *Response) Page() (*Page, error) {
	page := &Page{}

	if len(r.Query.Redirects) > 0 {
		page.Redirect = &r.Query.Redirects[0]
	}

	for _, p := range r.Query.Pages {
		url, err := url.QueryUnescape(p.Canonicalurl)
		if err != nil {
			url = p.Canonicalurl
		}
		page.ID = p.Pageid
		page.Title = p.Title
		page.Content = strings.Replace(p.Extract, "\n", "\n\n", -1)
		page.Language = p.Pagelanguage
		page.URL = url

		break
	}

	return page, nil
}

type query struct {
	Redirects []redirect
	Pages     map[string]page
}

type redirect struct {
	From string
	To   string
}

type page struct {
	Pageid       int
	Ns           int
	Title        string
	Extract      string
	Contentmodel string
	Pagelanguage string
	Touched      time.Time
	Fullurl      string
	Canonicalurl string
}
