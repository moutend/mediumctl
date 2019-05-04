package main

import (
	"log"
	"os"

	medium "github.com/moutend/go-medium"
	"github.com/spf13/cobra"
)

var (
	debug  bool
	client *medium.Client
	user   *medium.User
)

var rootCommand = &cobra.Command{
	Use: "mediumctl",
	PersistentPreRunE: func(c *cobra.Command, args []string) (err error) {
		switch c.Use {
		case "info", "user", "publication", "refresh":
			break
		default:
			return nil
		}

		token, err := readToken()
		if err != nil {
			return err
		}

		client = medium.NewClient(token.ApplicationID, token.ApplicationSecret, token.AccessToken)
		if debug {
			client.SetLogger(log.New(os.Stdout, "Debug: ", 0))
		}

		user, err = client.User()
		if err != nil {
			return err
		}

		return nil
	},
	PersistentPostRunE: func(c *cobra.Command, args []string) error {
		command := c
		commands := []string{}

		for {
			commands = append(commands, command.Use)

			if command = command.Parent(); command == nil {
				break
			}
		}

		return nil
	},
}

func init() {
	rootCommand.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug enable flag")
}
