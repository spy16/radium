package sources

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spy16/radium"
)

const learnXInYURL = "https://github.com/adambard/learnxinyminutes-docs/blob/master/%s.html.markdown"

// NewLearnXInYMins initializes a radium.Source implementation using
// LearnXInY github repository as the source of reference.
func NewLearnXInYMins() *LearnXInY {
	lxy := &LearnXInY{}
	return lxy
}

// LearnXInY implements radium.Source using LearnXInY github
// repository
type LearnXInY struct {
}

// Search attempts to download the appropriate markdown file from learn-x-in-y
// repository and format it as a result
func (lxy LearnXInY) Search(ctx context.Context, query radium.Query) ([]radium.Article, error) {
	var rs []radium.Article
	lang := strings.Replace(query.Text, " ", "-", -1)
	if res, err := lxy.getLanguageMarkdown(ctx, lang); err == nil {
		rs = append(rs, *res)
	}
	return rs, nil
}

func (lxy LearnXInY) getLanguageMarkdown(ctx context.Context, language string) (*radium.Article, error) {
	ghURL := fmt.Sprintf(learnXInYURL, url.QueryEscape(language))
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, _ := http.NewRequest(http.MethodGet, ghURL, nil)
	req.Header.Set("User-Agent", "curl/7.54.0")
	req.WithContext(ctx)

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
	result.Title = language
	result.Attribs = map[string]string{}
	result.License = "CC BY-SA 3.0"
	return result, nil
}
