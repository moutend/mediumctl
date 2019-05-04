package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "v0.4.1"
	commit  = "latest"
)

var versionCommand = &cobra.Command{
	Use:     "version",
	Short:   "Show the version of this command",
	Long:    "Show the version of this command",
	Aliases: []string{"v"},
	RunE: func(c *cobra.Command, args []string) (err error) {
		fmt.Printf("%s-%s\n", version, commit)

		return nil
	},
}

func init() {
	rootCommand.AddCommand(versionCommand)
}
