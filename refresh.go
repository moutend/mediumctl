package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var refreshCommand = &cobra.Command{
	Use:     "refresh",
	Short:   "Refresh the existing API token",
	Long:    "Refresh the existing API token",
	Aliases: []string{"r"},
	RunE: func(c *cobra.Command, args []string) (err error) {
		token, err := readToken()
		if err != nil {
			return err
		}

		newToken, err := client.RefreshToken(token.RefreshToken)
		if err != nil {
			return
		}
		if err := writeToken(token.ApplicationID, token.ApplicationSecret, newToken); err != nil {
			return err
		}

		fmt.Println("Done")

		return nil
	},
}

func init() {
	rootCommand.AddCommand(refreshCommand)
}
