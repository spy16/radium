package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/shivylp/radium"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

func writeOut(cmd *cobra.Command, v interface{}) {
	ugly, _ := cmd.Flags().GetBool("ugly")
	asJSON, _ := cmd.Flags().GetBool("json")
	if ugly {
		rawDump(v, asJSON)
	} else {
		tryPrettyPrint(v)
	}
}

func tryPrettyPrint(v interface{}) {
	if article, ok := v.(radium.Article); ok {
		fmt.Println(article.Content)
	} else {
		rawDump(v, true)
	}
}

func rawDump(v interface{}, asJSON bool) {
	var data []byte
	if asJSON {
		data, _ = json.MarshalIndent(v, "", "    ")
	} else {
		data, _ = yaml.Marshal(v)
	}
	os.Stdout.Write(data)
	os.Stdout.Write([]byte("\n"))
}
