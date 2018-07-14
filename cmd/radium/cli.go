package main

import (
	"fmt"
	"log"
	"reflect"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newCLI() *cobra.Command {
	cfg := &config{}
	rootCmd := newRootCmd(cfg)

	initConfig := func() {
		viper.SetConfigName("radium")
		viper.SetConfigType("yaml")

		viper.AddConfigPath("./")
		if hd, err := homedir.Dir(); err == nil {
			viper.AddConfigPath(hd)
		}
		viper.AutomaticEnv()
		viper.BindPFlags(rootCmd.PersistentFlags())

		viper.ReadInConfig()
		if err := viper.Unmarshal(cfg); err != nil {
			log.Fatalf("config err: %s\n", err)
		}
	}

	cobra.OnInitialize(initConfig)
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
		ins := getNewRadiumInstance()
		srcs := ins.GetSources()

		if L := len(srcs); L > 0 {
			fmt.Printf("%d source(s) available:\n", L)
			for name, src := range srcs {
				ty := reflect.TypeOf(src)
				fmt.Printf("* %s (Type: %s)\n", name, ty.String())
			}
		} else {
			fmt.Println("No sources configured")
		}
	}

	return cmd
}
