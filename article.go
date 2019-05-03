package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ericaro/frontmatter"
	medium "github.com/moutend/go-medium"
)

func parseArticle(filename string) (article medium.Article, publicationNumber int, err error) {
	type FrontmatterOption struct {
		Title             string   `fm:"title"`
		ContentFormat     string   `fm:"contentFormat"`
		Content           string   `fm:"content"`
		Curl              string   `fm:"curl"`
		Tags              []string `fm:"tags"`
		Status            string   `fm:"status"`
		PublishedAt       string   `fm:"publishedAt"`
		PublicationNumber int      `fm:"publicationNumber"`
		License           string   `fm:"license"`
		NotifyFollowers   bool     `fm:"notifyFollowers"`
	}

	var file []byte
	var fo *FrontmatterOption = &FrontmatterOption{}

	if file, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	if err = frontmatter.Unmarshal(file, fo); err != nil {
		return
	}
	if strings.HasSuffix(filename, "html") || strings.HasSuffix(filename, "htm") {
		fo.ContentFormat = "html"
	} else {
		fo.ContentFormat = "markdown"
	}

	article = medium.Article{
		Title:           fo.Title,
		ContentFormat:   fo.ContentFormat,
		Content:         fo.Content,
		CanonicalURL:    fo.Curl,
		Tags:            fo.Tags,
		PublishStatus:   fo.Status,
		PublishedAt:     fo.PublishedAt,
		License:         fo.License,
		NotifyFollowers: false,
	}
	publicationNumber = fo.PublicationNumber

	return
}

func showPostedArticleInfo(p *medium.PostedArticle) {
	fmt.Printf("Your article was successfully posted.\n\n")
	fmt.Printf("title: %s\n", p.Title)
	fmt.Printf("URL: %s\n", p.URL)

	if p.CanonicalURL != "" {
		fmt.Printf("canonicalURL: %s\n", p.CanonicalURL)
	}

	if p.PublishStatus == "" {
		fmt.Println("publishStatus: public")
	} else {
		fmt.Printf("publishStatus: %s\n", p.PublishStatus)
	}
	if len(p.Tags) > 0 {
		fmt.Printf("tags: %s\n", strings.Join(p.Tags, " "))
	}
}
