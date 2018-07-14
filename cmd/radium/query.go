package main

import (
	"context"
	"fmt"
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

	var tags []string
	cmd.Flags().StringSliceVarP(&tags, "tag", "t", []string{}, "Tags to narrow the search scope")

	cmd.Run = func(_ *cobra.Command, args []string) {
		query := radium.Query{}
		query.Text = args[0]
		query.Tags = map[string]string{}

		for _, tag := range tags {
			parts := strings.Split(tag, ":")
			if len(parts) == 2 {
				query.Tags[parts[0]] = parts[1]
			} else {
				fmt.Println("Err: invalid tag format. must be <name>:<value>")
				os.Exit(1)
			}
		}

		ctx := context.Background()
		ins := getNewRadiumInstance()
		rs, err := ins.Search(ctx, query)
		if err != nil {
			writeOut(cmd, map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			if len(rs) == 1 {
				writeOut(cmd, rs[0])
			} else {
				writeOut(cmd, rs)
			}
		}

	}

	return cmd
}
