package main

import (
	"github.com/shivylp/radium"
	"github.com/shivylp/radium/sources"
)

func getNewRadiumInstance() *radium.Instance {
	ins := radium.New(nil, nil)

	ins.RegisterSource("cheat.sh", sources.NewCheatSh())
	ins.RegisterSource("tldr", sources.NewTLDR())
	ins.RegisterSource("learnxinyminutes", sources.NewLearnXInYMins())

	return ins
}
