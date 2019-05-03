package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "v0.4.0"
	commit  = "latest"
)

var versionCommand = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	RunE: func(c *cobra.Command, args []string) (err error) {
		fmt.Printf("%s-%s\n", version, commit)

		return nil
	},
}

func init() {
	rootCommand.AddCommand(versionCommand)
}
