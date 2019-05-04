package main

import "github.com/spf13/cobra"

var userCommand = &cobra.Command{
	Use:     "user",
	Short:   "Publish the article to the user page",
	Long:    "Publish the article to the user page",
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
