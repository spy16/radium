package main

import (
	"strings"

	"github.com/shivylp/radium"
	"github.com/shivylp/radium/sources"
	"github.com/shivylp/radium/sources/cheatsh"
	"github.com/shivylp/radium/sources/duckduckgo"
	"github.com/shivylp/radium/sources/wikipedia"
)

func getNewRadiumInstance(cfg config) *radium.Instance {
	ins := radium.New(nil, nil)

	for _, src := range cfg.Sources {
		switch strings.ToLower(strings.TrimSpace(src)) {
		case "cheatsh", "cheat.sh":
			ins.RegisterSource("cheat.sh", cheatsh.New())
		case "learnxiny", "lxy", "learnxinyminutes":
			ins.RegisterSource("learnxinyminutes", sources.NewLearnXInYMins())
		case "tldr":
			ins.RegisterSource("tldr", sources.NewTLDR())
		case "wiki", "wikipedia":
			ins.RegisterSource("wikipedia", wikipedia.New("https://%s.wikipedia.org/w/api.php"))
		case "duckduckgo", "ddg":
			ins.RegisterSource("duckduckgo", duckduckgo.New())
		case "radium", "rad":
			ins.RegisterSource("radium", sources.NewRadium(cfg.RadiumURL))
		default:
			ins.Fatalf("unknown source type: %s", src)
		}
	}

	return ins
}
