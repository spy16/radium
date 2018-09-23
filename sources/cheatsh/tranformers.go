package cheatsh

import (
	"regexp"

	"github.com/shivylp/radium"
)

var (
	langQuery = regexp.MustCompile("(.*) in (\\w+)")
	languages = []string{
		"go",
		"golang",
		"java",
		"python",
		"lua",
		"c",
		"c++",
		"cpp",
		"ruby",
		"rails",
	}
)

func transformLanguageQuery(query *radium.Query) {
	matches := langQuery.FindAllStringSubmatch(query.Text, -1)
	if len(matches) >= 1 && contains(languages, matches[0][2]) {
		query.Text = matches[0][1]
		query.Attribs["language"] = matches[0][2]
	}
}

func contains(strs []string, str string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}

	return false
}
