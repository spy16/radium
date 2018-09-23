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
radium is a platform for reference content. radium lets
users publish, manage and view content from terminal and web.

radium can be used in multiple ways - as an offline, personal
wiki; a server to serve content; a clipboard service to monitor
and replace queries with solutions as they come in!
`,
	}

	rootCmd.PersistentFlags().BoolP("ugly", "u", false, "Print raw output as yaml or json")
	rootCmd.PersistentFlags().Bool("json", false, "Print output as JSON")
	rootCmd.PersistentFlags().StringSlice("sources", nil, "Enable sources")

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
				ty := reflect.TypeOf(src)
				fmt.Printf("%d. %s (Type: %s)\n", order+1, src.Name, ty.String())
			}
		} else {
			fmt.Println("No sources configured")
		}
	}

	return cmd
}
