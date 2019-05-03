// Author: Yoshiyuki Koyanagi <moutend@gmail.com>
// License: mIT

// Package main implements mediumctl.
package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/ericaro/frontmatter"
	medium "github.com/moutend/go-medium"
)

type Token struct {
	ApplicationID     string
	ApplicationSecret string
	AccessToken       string
	ExpiresAt         int
	RefreshToken      string
}

var (
	version  = "v0.3.0"
	revision = "dev"
)

func showPostedArticleInfo(p *medium.PostedArticle) {
	fmt.Printf("Your article was successfully posted.\n\n")
	fmt.Printf("title: %s\n", p.Title)
	if p.PublishStatus == "" {
		fmt.Println("publishStatus: public")
	} else {
		fmt.Printf("publishStatus: %s\n", p.PublishStatus)
	}
	if len(p.Tags) > 0 {
		fmt.Printf("tags: %s\n", strings.Join(p.Tags, " "))
	}
	fmt.Printf("URL: %s\n", p.URL)
	if p.CanonicalURL != "" {
		fmt.Printf("canonicalURL: %s\n", p.CanonicalURL)
	}
	return
}

func parseArticle(filename string) (article medium.Article, publicationNumber int, err error) {
	type FrontmatterOption struct {
		Title             string   `fm:"title"`
		ContentFormat     string   `fm:"contentFormat"`
		Content           string   `fm:"content"`
		CanonicalURL      string   `fm:"canonicalUrl"`
		Tags              []string `fm:"tags"`
		PublishStatus     string   `fm:"publishStatus"`
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
		CanonicalURL:    fo.CanonicalURL,
		Tags:            fo.Tags,
		PublishStatus:   fo.PublishStatus,
		PublishedAt:     fo.PublishedAt,
		License:         fo.License,
		NotifyFollowers: false,
	}
	publicationNumber = fo.PublicationNumber

	return
}

func getUser(token *Token, debugFlag bool) (u *medium.User, err error) {
	c := medium.NewClient(token.ApplicationID, token.ApplicationSecret, token.AccessToken)
	if debugFlag {
		c.SetLogger(log.New(os.Stdout, "debug: ", 0))
	}
	return c.User()
}

func getCode(clientID string, redirectURL *url.URL) (code string, err error) {
	listener, err := net.Listen("tcp", redirectURL.Hostname()+":"+redirectURL.Port())
	if err != nil {
		return
	}
	defer listener.Close()

	responseChann := make(chan string)

	go http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(`<script>window.open("about:blank","_self").close()</script>`))
		w.(http.Flusher).Flush()

		e := req.FormValue("error")
		if e != "" {
			err = fmt.Errorf(e)
		}
		responseChann <- req.FormValue("code")
	}))

	stateBytes := make([]byte, 88)
	_, err = rand.Read(stateBytes)
	if err != nil {
		return
	}
	state := fmt.Sprintf("%x", stateBytes)
	scope := "basicProfile,listPublications,publishPost"
	query := fmt.Sprintf("client_id=%s&scope=%s&state=%s&response_type=code&redirect_uri=%s", clientID, scope, state, redirectURL)
	uri := "https://medium.com/m/oauth/authorize?" + query
	fmt.Println("Please open this URL:", uri)

	select {
	case code = <-responseChann:
		break
	case <-time.After(60 * time.Second):
		err = fmt.Errorf("timeout")
		break
	}
	return
}

func getTokenFilePath() (tokenFilePath string, err error) {
	var u *user.User
	const tokenFileName = ".mediumctl"

	if u, err = user.Current(); err != nil {
		return
	}
	tokenFilePath = filepath.Join(u.HomeDir, tokenFileName)
	return
}

func writeToken(clientID, clientSecret string, t *medium.Token) (err error) {
	var file []byte
	var tokenFilePath string

	file, err = json.Marshal(Token{
		ApplicationID:     clientID,
		ApplicationSecret: clientSecret,
		AccessToken:       t.AccessToken,
		ExpiresAt:         t.ExpiresAt,
		RefreshToken:      t.RefreshToken,
	})
	if err != nil {
		return
	}
	if tokenFilePath, err = getTokenFilePath(); err != nil {
		return
	}
	err = ioutil.WriteFile(tokenFilePath, file, 0644)
	return
}

func readToken() (token *Token, err error) {
	var file []byte
	var tokenFilePath string

	if tokenFilePath, err = getTokenFilePath(); err != nil {
		return
	}
	if file, err = ioutil.ReadFile(tokenFilePath); err != nil {
		err = fmt.Errorf("API token is not set")
		return
	}

	token = &Token{}
	if err = json.Unmarshal(file, token); err != nil {
		return
	}
	return
}

func main() {
	var err error

	if err = run(os.Args); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func run(args []string) (err error) {
	var token *Token

	if len(args) < 2 {
		return helpCommand(args)
	}
	switch args[1] {
	case "v":
		err = versionCommand(args)
	case "version":
		err = versionCommand(args)
	case "h":
		err = helpCommand(args)
	case "help":
		err = helpCommand(args)
	case "o":
		err = authCommand(args)
	case "oauth":
		err = authCommand(args)
	}
	if token, err = readToken(); err != nil {
		return
	}
	switch args[1] {
	case "r":
		err = refreshCommand(token, args)
	case "refresh":
		err = refreshCommand(token, args)
	case "i":
		err = infoCommand(token, args)
	case "info":
		err = infoCommand(token, args)
	case "p":
		err = publicationCommand(token, args)
	case "publication":
		err = publicationCommand(token, args)
	case "u":
		err = userCommand(token, args)
	case "user":
		err = userCommand(token, args)
	default:
		fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s subcommand.\n", args[0], args[1], args[0])
		err = helpCommand(args)
	}
	return
}

func authCommand(args []string) (err error) {
	var clientIDFlag string
	var clientSecretFlag string
	var debugFlag bool
	var redirectURLFlag string
	var redirectURL *url.URL
	var code string
	var token *medium.Token

	f := flag.NewFlagSet(fmt.Sprintf("%s %s", args[0], args[1]), flag.ExitOnError)
	f.StringVar(&redirectURLFlag, "u", "", "Redirect URL for OAuth application.")
	f.StringVar(&clientIDFlag, "i", "", "Client ID of OAuth application.")
	f.StringVar(&clientSecretFlag, "s", "", "Client secret of OAuth application.")
	f.BoolVar(&debugFlag, "debug", false, "Enable debug output.")
	f.Parse(args[2:])

	if redirectURLFlag == "" {
		return fmt.Errorf("please specify redirect URL")
	}
	if clientIDFlag == "" {
		return fmt.Errorf("please specify client ID")
	}
	if clientSecretFlag == "" {
		return fmt.Errorf("please specify client secret")
	}
	if redirectURL, err = url.Parse(redirectURLFlag); err != nil {
		return
	}
	if code, err = getCode(clientIDFlag, redirectURL); err != nil {
		return
	}

	c := medium.NewClient(clientIDFlag, clientSecretFlag, "")
	if debugFlag {
		c.SetLogger(log.New(os.Stdout, "debug: ", 0))
	}
	if token, err = c.Token(code, redirectURLFlag); err != nil {
		return
	}
	if err = writeToken(clientIDFlag, clientSecretFlag, token); err != nil {
		return
	}

	fmt.Println("Your API token was successfully saved ")

	return
}

func refreshCommand(token *Token, args []string) (err error) {
	var debugFlag bool
	var refreshedToken *medium.Token

	f := flag.NewFlagSet(fmt.Sprintf("%s %s", args[0], args[1]), flag.ExitOnError)
	f.BoolVar(&debugFlag, "debug", false, "Enable debug output.")
	f.Parse(args[2:])

	c := medium.NewClient(token.ApplicationID, token.ApplicationSecret, "")
	if debugFlag {
		c.SetLogger(log.New(os.Stdout, "debug: ", 0))
	}
	if refreshedToken, err = c.RefreshToken(token.RefreshToken); err != nil {
		return
	}
	if err = writeToken(token.ApplicationID, token.ApplicationSecret, refreshedToken); err != nil {
		return
	}
	fmt.Println("Your API token has been refreshed.")
	return
}

func infoCommand(token *Token, args []string) (err error) {
	var debugFlag bool
	var u *medium.User
	var ps []*medium.Publication

	f := flag.NewFlagSet(fmt.Sprintf("%s %s", args[0], args[1]), flag.ExitOnError)
	f.BoolVar(&debugFlag, "debug", false, "Enable debug output.")
	f.Parse(args[2:])

	if u, err = getUser(token, debugFlag); err != nil {
		return
	}

	fmt.Printf("You are logged in as:\n\n")
	fmt.Printf("Name: %s\n", u.Name)
	fmt.Printf("Username: %s\n", u.Username)
	fmt.Printf("URL: %s", u.URL)
	fmt.Printf("\n")

	if ps, err = u.Publications(); err != nil {
		return
	}
	if len(ps) == 0 {
		fmt.Println("You have no publications yet.")
		return
	}

	fmt.Printf("\nYou have publication(s) below:\n\n")
	for i, p := range ps {
		fmt.Printf("Number: %d\n", i)
		fmt.Printf("Name: %s\n", p.Name)
		fmt.Printf("Description: %s\n", p.Description)
		fmt.Printf("URL: %s\n\n", p.URL)
	}
	return
}

func userCommand(token *Token, args []string) (err error) {
	var debugFlag bool
	var article medium.Article
	var u *medium.User
	var postedArticle *medium.PostedArticle

	f := flag.NewFlagSet(fmt.Sprintf("%s %s", args[0], args[1]), flag.ExitOnError)
	f.BoolVar(&debugFlag, "debug", false, "Enable debug output.")
	f.Parse(args[2:])

	if article, _, err = parseArticle(f.Args()[0]); err != nil {
		return
	}
	if u, err = getUser(token, debugFlag); err != nil {
		return
	}
	if postedArticle, err = u.Post(article); err != nil {
		return
	}
	showPostedArticleInfo(postedArticle)

	return
}

func publicationCommand(token *Token, args []string) (err error) {
	var debugFlag bool
	var article medium.Article
	var publicationNumber int
	var u *medium.User
	var ps []*medium.Publication
	var postedArticle *medium.PostedArticle

	f := flag.NewFlagSet(fmt.Sprintf("%s %s", args[0], args[1]), flag.ExitOnError)
	f.BoolVar(&debugFlag, "debug", false, "Enable debug output.")
	f.Parse(args[2:])

	if article, publicationNumber, err = parseArticle(f.Args()[0]); err != nil {
		return
	}
	if u, err = getUser(token, debugFlag); err != nil {
		return
	}
	if ps, err = u.Publications(); err != nil {
		return
	}
	if len(ps) == 0 {
		err = fmt.Errorf("you have no publications yet")
		return
	}
	if publicationNumber < 0 || publicationNumber > len(ps)-1 {
		err = fmt.Errorf("publication number '%d' is invalid", publicationNumber)
		return
	}
	if postedArticle, err = ps[publicationNumber].Post(article); err != nil {
		return
	}
	showPostedArticleInfo(postedArticle)

	return
}

func versionCommand(args []string) (err error) {
	fmt.Printf("%s-%s\n", version, revision)
	return
}

func helpCommand(args []string) (err error) {
	fmt.Println(`usage: mediumctl <command> [options]

Commands:
  oauth, o
    Setting up API token with OAuth.
  refresh, r
    Refresh existing API token.
  info, i
    Show the information about current user and its publications.
  user, u
    Post HTML or Markdown file to current user profile.
  publication, p
    Post HTML or Markdown file to current user's publication.
  version, v
  Show version and revision information.
  help, h
    Show this message.

For more information, please see https://github.com/moutend/mediumctl.`)
	return
}
