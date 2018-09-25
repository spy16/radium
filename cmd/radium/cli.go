package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
)

func newCLI() *cobra.Command {
	cfg := &config{}
	rootCmd := newRootCmd(cfg)

	cobra.OnInitialize(func() {
		initConfig(cfg, rootCmd)
	})
	return rootCmd
}

func newRootCmd(cfg *config) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "radium <command>",
		Short: "radium is a platform for reference content",
		Long: `
radium is a platform (client and optional server) for viewing reference
articles, cheat sheets etc. right from a shell.

radium can be used in multiple ways - as an offline, wiki from CLI,
a server to serve content, a clipboard service to monitor and replace
queries with solutions as they come in!

Examples:
  // general query which sequentially executes configured sources
  // until the first result is found
  $ radium query "append file in go"

  // custom query which looks only 'wikipedia' source and searches
  // for the query
  $ radium query "apple" -s wikipedia

  // custom query which executes all configured sources concurrently
  // and returns all the results
  $ radium query "apple" --strategy concurrent
`,
	}

	rootCmd.PersistentFlags().BoolP("ugly", "u", false, "Print raw output as yaml or json")
	rootCmd.PersistentFlags().BoolP("json", "j", false, "Print output as JSON")
	rootCmd.PersistentFlags().StringSliceP("sources", "s", nil, "Enable sources")
	rootCmd.PersistentFlags().StringP("strategy", "S", "1st", "Default strategy to use")

	rootCmd.AddCommand(newServeCmd(cfg))
	rootCmd.AddCommand(newQueryCmd(cfg))
	rootCmd.AddCommand(newListSources(cfg))

	return rootCmd
}

func newListSources(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sources",
		Short: "List registered sources",
		Args:  cobra.ExactArgs(0),
	}

	cmd.Run = func(_ *cobra.Command, args []string) {
		ins := getNewRadiumInstance(*cfg)
		srcs := ins.GetSources()

		if L := len(srcs); L > 0 {
			fmt.Printf("%d source(s) available:\n", L)
			fmt.Printf("%s\n", strings.Repeat("-", 20))
			for order, src := range srcs {
				ty := reflect.TypeOf(src.Source)
				fmt.Printf("%d. %s (Type: %s)\n", order+1, src.Name, ty.String())
			}
		} else {
			fmt.Println("No sources configured")
		}
	}

	return cmd
}
