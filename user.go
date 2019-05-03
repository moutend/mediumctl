package main

import "github.com/spf13/cobra"

var userCommand = &cobra.Command{
	Use:     "user",
	Aliases: []string{"u"},
	RunE: func(c *cobra.Command, args []string) (err error) {
		if len(args) == 0 {
			return nil
		}

		article, _, err := parseArticle(args[0])
		if err != nil {
			return err
		}

		postedArticle, err := user.Post(article)
		if err != nil {
			return err
		}

		showPostedArticleInfo(postedArticle)

		return nil
	},
}

func init() {
	rootCommand.AddCommand(userCommand)
}
