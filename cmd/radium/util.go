package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spy16/radium"
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
	switch v.(type) {
	case radium.Article, *radium.Article:
		fmt.Println((v.(radium.Article)).Content)
	case []radium.Article:
		results := v.([]radium.Article)
		if len(results) == 1 {
			tryPrettyPrint(results[0])
		} else {
			for num, res := range results {
				fmt.Printf("%d. %s [source=%s]\n", num+1, strings.TrimSpace(res.Title), res.Source)
				fmt.Println(strings.Repeat("-", 20))
				fmt.Println(res.Content)
				fmt.Println(strings.Repeat("=", 20))
			}
		}
	case error:
		fmt.Printf("error: %s\n", v)
	default:
		rawDump(v, false)
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
