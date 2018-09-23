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
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Request sets up the request against the api with the correct parameters
// and has functionality to fetch the data and convert it to a response.
type Request struct {
	*url.URL
}

// NewRequest creates a new request against baseURL for language.
// Language is interpolated in the baseURL if asked, if not it is ignored.
// Query is the title of the page to fetch.
// Returns an error if the URL can not be parsed.
func NewRequest(baseURL, query, language string) (*Request, error) {
	if strings.Contains(baseURL, "%s") {
		baseURL = fmt.Sprintf(baseURL, language)
	}
	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	v := url.Query()
	v.Set("action", "query")
	v.Set("prop", "extracts|info")
	v.Set("format", "json")
	v.Set("exintro", "")
	v.Set("explaintext", "")
	v.Set("inprop", "url")
	v.Set("redirects", "")
	v.Set("converttitles", "")
	v.Set("titles", query)
	url.RawQuery = v.Encode()

	return &Request{url}, nil
}

// Execute fetches the data and decodes it into a Response.
// Returns an error if the data could not be retrived or the decoding fails.
func (r *Request) Execute(ctx context.Context, noCheckCert bool) (*Response, error) {
	client := &http.Client{}

	if noCheckCert {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	}

	req, err := http.NewRequest(http.MethodGet, r.String(), nil)
	if err != nil {
		return nil, err
	}
	req.WithContext(ctx)

	data, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	d := json.NewDecoder(data.Body)
	resp := &Response{}
	err = d.Decode(resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
