package sources

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/shivylp/radium"
)

// NewTLDR initializes a radium.Source implementation using
// the popular tldr-pages as a source
func NewTLDR() *TLDR {
	tldr := &TLDR{}
	return tldr
}

// TLDR implements radium.Source using tldr-pages
// as the source of reference
type TLDR struct {
}

// Search for a particular query in tldr-pages repository
func (tldr TLDR) Search(query radium.Query) ([]radium.Article, error) {
	var rs []radium.Article

	tool := strings.Replace(query.Text, " ", "-", -1)
	platform := "common"

	if val, found := query.Tags["platform"]; found {
		platform = val
	}

	res, err := tldr.getPlatformToolInfo(tool, platform)
	if err == nil {
		rs = append(rs, *res)
	}
	return rs, nil
}

func (tldr TLDR) getPlatformToolInfo(tool, platform string) (*radium.Article, error) {
	rawGitURL := "https://raw.githubusercontent.com/tldr-pages/tldr/master/pages/%s/%s.md"

	ghURL := fmt.Sprintf(rawGitURL, url.QueryEscape(platform), url.QueryEscape(tool))
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, _ := http.NewRequest(http.MethodGet, ghURL, nil)
	req.Header.Set("User-Agent", "curl/7.54.0")

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

	result := &radium.Article{}
	result.Content = string(data)
	result.ContentType = "markdown"
	result.Title = tool
	result.Tags = map[string]string{
		"platform": platform,
	}
	result.License = "The MIT License (MIT)"
	return result, nil
}
