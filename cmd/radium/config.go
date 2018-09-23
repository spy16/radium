package main

import (
	"log"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initConfig(cfg *config, rootCmd *cobra.Command) {
	viper.SetDefault("sources", "cheatsh,learnxiny,tldr")
	viper.SetDefault("radiumservers", []string{})
	viper.RegisterAlias("radiumservers", "radium_servers")

	viper.SetConfigName("radium")
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

// config struct is used to store CLI configurations. configuration
// values are read into this struct using viper
type config struct {
	Sources       []string `json:"sources"`
	RadiumServers []string `json:"radium_servers"`
}
