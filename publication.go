package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var publicationCommand = &cobra.Command{
	Use:     "publication",
	Short:   "Publish the article to the publication page",
	Long:    "Publish the article to the publication page",
	Aliases: []string{"p"},
	RunE: func(c *cobra.Command, args []string) (err error) {
		if len(args) == 0 {
			return nil
		}

		article, publicationNumber, err := parseArticle(args[0])
		if err != nil {
			return err
		}

		publications, err := user.Publications()
		if err != nil {
			return
		}
		if len(publications) == 0 {
			return fmt.Errorf("publications not found")
		}
		if publicationNumber < 0 || publicationNumber > len(publications)-1 {
			return fmt.Errorf("publication number '%d' is invalid", publicationNumber)
		}

		postedArticle, err := publications[publicationNumber].Post(article)
		if err != nil {
			return err
		}

		showPostedArticleInfo(postedArticle)

		return nil
	},
}

func init() {
	rootCommand.AddCommand(publicationCommand)
}
