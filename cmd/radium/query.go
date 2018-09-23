package main

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/shivylp/radium"
	"github.com/spf13/cobra"
)

func newQueryCmd(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query <query>",
		Short:   "Run a query",
		Aliases: []string{"q", "search"},
		Args:    cobra.ExactArgs(1),
	}

	var attribs []string
	var strategy string
	cmd.Flags().StringSliceVarP(&attribs, "attr", "a", []string{}, "Attributes to narrow the search scope")
	cmd.Flags().StringVarP(&strategy, "strategy", "s", "concurrent", "Strategy to use for executing sources")

	cmd.Run = func(_ *cobra.Command, args []string) {
		query := radium.Query{}
		query.Text = args[0]
		query.Attribs = map[string]string{}

		for _, attrib := range attribs {
			parts := strings.Split(attrib, ":")
			if len(parts) == 2 {
				query.Attribs[parts[0]] = parts[1]
			} else {
				writeOut(cmd, errors.New("invalid attrib format. must be <name>:<value>"))
				os.Exit(1)
			}
		}

		ctx := context.Background()
		ins := getNewRadiumInstance(*cfg)
		rs, err := ins.Search(ctx, query, strategy)
		if err != nil {
			writeOut(cmd, err)
			os.Exit(1)
		}

		if len(rs) == 1 {
			writeOut(cmd, rs[0])
		} else {
			writeOut(cmd, rs)
		}
	}

	return cmd
}
