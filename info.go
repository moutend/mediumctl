package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var infoCommand = &cobra.Command{
	Use:     "info",
	Short:   "Show the info about yourself",
	Long:    "Show the info about yourself",
	Aliases: []string{"i"},
	RunE: func(c *cobra.Command, args []string) (err error) {
		fmt.Printf("Name: %s\n", user.Name)
		fmt.Printf("Username: %s\n", user.Username)
		fmt.Printf("URL: %s", user.URL)

		publications, err := user.Publications()
		if err != nil {
			return err
		}
		if len(publications) == 0 {
			fmt.Println("---")
			fmt.Println("You have no publications.")

			return nil
		}

		for i, p := range publications {
			fmt.Println("---")
			fmt.Printf("Number: %d\n", i)
			fmt.Printf("Name: %s\n", p.Name)
			fmt.Printf("Description: %s\n", p.Description)
			fmt.Printf("URL: %s\n\n", p.URL)
		}

		return nil
	},
}

func init() {
	rootCommand.AddCommand(infoCommand)
}
