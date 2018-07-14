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

// NewCheatSh initializes a radium.Source implementation
// using http://cheat.sh as the source
func NewCheatSh() *CheatSh {
	csh := &CheatSh{}
	return csh
}

// CheatSh implements radium.Source interface using
// cheat.sh as the source
type CheatSh struct {
}

// Search performs an HTTP request to http://cheat.sh to find
// results matching the given query.
func (csh CheatSh) Search(query radium.Query) ([]radium.Article, error) {
	var results []radium.Article

	if lang, found := query.Tags["language"]; found {
		color := false
		if val, found := query.Tags["color"]; found {
			if val == "yes" || val == "true" {
				color = true
			}
		}
		res, err := csh.makeLangRequest(query.Text, lang, color)
		if err == nil {
			results = append(results, *res)
		}
	}
	return results, nil
}

func (csh CheatSh) makeLangRequest(q string, lang string, color bool) (*radium.Article, error) {
	queryStr := url.QueryEscape(strings.Replace(q, " ", "+", -1))
	csURL := fmt.Sprintf("http://cheat.sh/%s/%s", url.QueryEscape(lang), queryStr)

	if !color {
		csURL += "?T"
	}

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, _ := http.NewRequest(http.MethodGet, csURL, nil)
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
	result.ContentType = "plaintext"
	result.Title = q
	result.Tags = map[string]string{
		"language": lang,
	}
	return result, nil
}
