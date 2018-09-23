package main

import "fmt"

var (
	version = "n/a"
	commit  = "n/a"
)

func main() {
	cli := newCLI()
	cli.Version = fmt.Sprintf("%s (commit: %s)", version, commit)
	cli.Execute()
}
