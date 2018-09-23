package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

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
	var deadline time.Duration
	cmd.Flags().StringSliceVarP(&attribs, "attr", "a", []string{}, "Attributes to narrow the search scope")
	cmd.Flags().DurationVarP(&deadline, "timeout", "t", 10*time.Second, "Timeout to use across all sources")

	cmd.Run = func(_ *cobra.Command, args []string) {
		query := radium.Query{}
		query.Text = args[0]
		query.Attribs = map[string]string{}

		for _, attrib := range attribs {
			parts := strings.Split(attrib, ":")
			if len(parts) == 2 {
				query.Attribs[parts[0]] = parts[1]
			} else if len(parts) == 1 {
				query.Attribs[parts[0]] = "true"
			} else {
				writeOut(cmd, errors.New("invalid attrib format. must be <name>[:<value>]"))
				os.Exit(1)
			}
		}

		strategy, err := cmd.Flags().GetString("strategy")
		if err != nil {
			strategy = "1st"
		}

		ctx := context.Background()
		var cancel func()
		if deadline.Seconds() > 0 {
			ctx, cancel = context.WithTimeout(ctx, deadline)
			defer cancel()
		}

		ins := getNewRadiumInstance(*cfg)
		rs, err := ins.Search(ctx, query, strategy)
		if err != nil {
			writeOut(cmd, err)
			os.Exit(1)
		}

		if len(rs) == 0 {
			fmt.Println("No results found")
		} else {
			writeOut(cmd, rs)
		}
	}

	return cmd
}
