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
	"strconv"
	"strings"
	"time"

	medium "github.com/moutend/go-medium"
	"github.com/skratchdot/open-golang/open"
)

type token struct {
	ApplicationID     string
	ApplicationSecret string
	AccessToken       string
	ExpiresAt         int
}

var (
	version       = "v0.1.1"
	revision      = "latest"
	tokenFilePath string
)

const tokenFileName = ".mediumctl"

func showPostedArticleInfo(p *medium.PostedArticle) {
	fmt.Println("Your article was successfully posted.")
	fmt.Printf("Title: %s\n", p.Title)
	fmt.Printf("Status: %s\n", p.PublishStatus)
	if len(p.Tags) > 0 {
		fmt.Printf("Tags: %s\n", strings.Join(p.Tags, " "))
	}
	fmt.Printf("URL: %s\n", p.URL)
	if p.CanonicalURL != "" {
		fmt.Printf("Canonical URL: %s\n", p.CanonicalURL)
	}
	return
}

func parseArticle(filename string) (article medium.Article, publicationNumber int, err error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	if len(b) == 0 {
		err = fmt.Errorf("%s is empty", filename)
		return
	}

	var (
		title        string
		tags         []string
		content      string
		format       string
		license      string
		status       string
		canonicalURL string
		notify       bool
	)
	format = "markdown"
	if strings.HasSuffix(filename, "html") || strings.HasSuffix(filename, "htm") {
		format = "html"
	}
	title = "untitled"
	status = "public"
	lines := strings.Split(string(b), "\n")

	for i, line := range lines[1:] {
		if strings.HasPrefix(line, "---") {
			content = strings.Join(lines[i+2:], "\n")
			break
		}
		if strings.HasPrefix(line, "number: ") {
			publicationNumber, err = strconv.Atoi(line[len("number: "):])
			if err != nil {
				return
			}
		}
		if strings.HasPrefix(line, "title: ") {
			title = line[len("title: "):]
		}
		if strings.HasPrefix(line, "tags: ") {
			tags = strings.Split(line[len("tags: "):], " ")
		}
		if strings.HasPrefix(line, "notify: true") {
			notify = true
		}
		if strings.HasPrefix(line, "status: ") {
			status = line[len("status: "):]
		}
		if strings.HasPrefix(line, "license: ") {
			license = line[len("license: "):]
		}
		if strings.HasPrefix(line, "canonicalURL: ") {
			canonicalURL = line[len("canonicalURL: "):]
		}
	}
	if content == "" {
		content = strings.Join(lines, "\n")
	}
	article = medium.Article{
		Title:           title,
		ContentFormat:   format,
		Content:         content,
		CanonicalURL:    canonicalURL,
		Tags:            tags,
		PublishStatus:   status,
		License:         license,
		NotifyFollowers: notify,
	}
	return
}
func getCode(clientID, redirectURI string) (code string, err error) {
	l, err := net.Listen("tcp", "192.168.1.107:4000")
	if err != nil {
		return
	}
	defer l.Close()

	type value struct {
		code  string
		error error
	}
	quit := make(chan value)
	go http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(`<script>window.open("about:blank","_self").close()</script>`))
		w.(http.Flusher).Flush()
		c := req.FormValue("code")
		e := req.FormValue("error")
		v := value{
			code:  c,
			error: nil,
		}
		if e != "" {
			v.error = fmt.Errorf(e)
		}
		quit <- v
	}))
	stateBytes := make([]byte, 88)
	_, err = rand.Read(stateBytes)
	if err != nil {
		return
	}
	state := fmt.Sprintf("%x", stateBytes)
	scope := "basicProfile,listPublications,publishPost"
	redirectURI = url.QueryEscape(redirectURI)
	q := fmt.Sprintf("client_id=%s&scope=%s&state=%s&response_type=code&redirect_uri=%s", clientID, scope, state, redirectURI)
	p := "https://medium.com/m/oauth/authorize?" + q
	if err = open.Start(p); err != nil {
		return
	}
	select {
	case v := <-quit:
		if v.error != nil {
			return "", v.error
		}
		return v.code, nil
	case <-time.After(60 * time.Second):
		return "", fmt.Errorf("timeout")
	}
}

func saveToken(clientID, clientSecret string, t *medium.Token) (err error) {
	b, err := json.Marshal(token{
		ApplicationID:     clientID,
		ApplicationSecret: clientSecret,
		AccessToken:       t.AccessToken,
		ExpiresAt:         t.ExpiresAt,
	})
	if err != nil {
		return
	}
	err = ioutil.WriteFile(tokenFilePath, b, 0644)
	return
}

func loadToken() (*token, error) {
	b, err := ioutil.ReadFile(tokenFilePath)
	if err != nil {
		return nil, fmt.Errorf("API token is not set. Please run 'auth' at first")
	}
	var t token
	err = json.Unmarshal(b, &t)
	return &t, err
}

func main() {
	err := run(os.Args)

	if err != nil {
		log.New(os.Stderr, "error: ", 0).Fatal(err)
		os.Exit(1)
	}

	os.Exit(0)
}
func run(args []string) (err error) {
	if len(args) < 2 {
		return helpCommand(args)
	}
	u, err := user.Current()
	if err != nil {
		return
	}
	tokenFilePath = filepath.Join(u.HomeDir, tokenFileName)
	switch args[1] {
	case "oauth":
		err = authCommand(args)
	case "o":
		err = authCommand(args)
	case "i":
		err = infoCommand(args)
	case "info":
		err = infoCommand(args)
	case "p":
		err = postCommand(args, false)
	case "publication":
		err = postCommand(args, false)
	case "u":
		err = postCommand(args, true)
	case "user":
		err = postCommand(args, true)
	case "v":
		err = versionCommand(args)
	case "version":
		err = versionCommand(args)
	case "h":
		err = helpCommand(args)
	case "help":
		err = helpCommand(args)
	default:
		fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s subcommand.\n", args[0], args[1], args[0])
		err = helpCommand(args)
	}
	return
}

func authCommand(args []string) (err error) {
	var (
		clientIDFlag     string
		clientSecretFlag string
		debugFlag        bool
		redirectURIFlag  string
	)

	f := flag.NewFlagSet(fmt.Sprintf("%s %s", args[0], args[1]), flag.ExitOnError)
	f.StringVar(&redirectURIFlag, "u", "", "Redirect URI for OAuth application.")
	f.StringVar(&clientIDFlag, "i", "", "Client ID of OAuth application.")
	f.StringVar(&clientSecretFlag, "s", "", "Client secret of OAuth application.")
	f.BoolVar(&debugFlag, "debug", false, "Enable debug output.")
	f.Parse(args[2:])
	if redirectURIFlag == "" {
		return fmt.Errorf("please specify redirect URI")
	}
	if clientIDFlag == "" {
		return fmt.Errorf("please specify client ID")
	}
	if clientSecretFlag == "" {
		return fmt.Errorf("please specify client secret")
	}

	code, err := getCode(clientIDFlag, redirectURIFlag)
	if err != nil {
		return
	}
	c := medium.NewClient(clientIDFlag, clientSecretFlag, "")
	if debugFlag {
		c.SetLogger(log.New(os.Stdout, "debug: ", 0))
	}
	token, err := c.Token(code, redirectURIFlag)
	if err != nil {
		return
	}
	if err = saveToken(clientIDFlag, clientSecretFlag, token); err != nil {
		return
	}
	fmt.Printf("Your API token was successfully saved in '%s'.\n", tokenFilePath)
	fmt.Println("Note: This file should be treated as the password and please do NOT expose it.")
	return
}

func infoCommand(args []string) (err error) {
	var (
		debugFlag bool
	)

	f := flag.NewFlagSet(fmt.Sprintf("%s %s", args[0], args[1]), flag.ExitOnError)
	f.BoolVar(&debugFlag, "debug", false, "Enable debug output.")
	f.Parse(args[2:])

	t, err := loadToken()
	if err != nil {
		return
	}
	c := medium.NewClient(t.ApplicationID, t.ApplicationSecret, t.AccessToken)
	if debugFlag {
		c.SetLogger(log.New(os.Stdout, "debug: ", 0))
	}
	u, err := c.User()
	if err != nil {
		return
	}
	fmt.Printf("You are logged in as:\n\n")
	fmt.Printf("Name: %s\n", u.Name)
	fmt.Printf("Username: %s\n", u.Username)
	fmt.Printf("URL: %s", u.URL)
	fmt.Printf("\n")

	ps, err := u.Publications()
	if err != nil {
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

func postCommand(args []string, userFlag bool) (err error) {
	var (
		debugFlag bool
	)

	f := flag.NewFlagSet(fmt.Sprintf("%s %s", args[0], args[1]), flag.ExitOnError)
	f.BoolVar(&debugFlag, "debug", false, "Enable debug output.")
	f.Parse(args[2:])

	article, publicationNumber, err := parseArticle(f.Args()[0])
	if err != nil {
		return
	}
	t, err := loadToken()
	if err != nil {
		return
	}
	c := medium.NewClient(t.ApplicationID, t.ApplicationSecret, t.AccessToken)
	if debugFlag {
		c.SetLogger(log.New(os.Stdout, "debug: ", 0))
	}
	u, err := c.User()
	if err != nil {
		return
	}
	if userFlag {
		p, err := u.Post(article)
		if err != nil {
			return err
		}
		showPostedArticleInfo(p)
		return nil
	}
	ps, err := u.Publications()
	if err != nil {
		return
	}
	if len(ps) == 0 {
		return fmt.Errorf("you have no publications yet")
	}
	if publicationNumber < 0 || publicationNumber > len(ps)-1 {
		err = fmt.Errorf("publication number '%d' is invalid", publicationNumber)
		return
	}
	p, err := ps[publicationNumber].Post(article)
	if err != nil {
		return
	}
	showPostedArticleInfo(p)
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
    Setting up API token for Medium with OAuth.
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
