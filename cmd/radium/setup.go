package main

import (
	"strings"

	"github.com/shivylp/radium"
	"github.com/shivylp/radium/sources"
)

func getNewRadiumInstance(cfg config) *radium.Instance {
	ins := radium.New(nil, nil)

	for _, src := range cfg.Sources {
		switch strings.ToLower(strings.TrimSpace(src)) {
		case "cheatsh", "cheat.sh":
			ins.RegisterSource("cheat.sh", sources.NewCheatSh())
		case "learnxiny", "lxy", "learnxinyminutes":
			ins.RegisterSource("learnxinyminutes", sources.NewLearnXInYMins())
		case "tldr":
			ins.RegisterSource("tldr", sources.NewTLDR())
		case "wiki", "wikipedia":
			ins.RegisterSource("wikipedia", sources.NewWikipedia())
		default:
			ins.Fatalf("unknown source type: %s", src)
		}
	}

	return ins
}
