// Author: Yoshiyuki Koyanagi <moutend@gmail.com>
// License: mIT

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

	medium "github.com/moutend/go-medium"
	"github.com/skratchdot/open-golang/open"
)

type Token struct {
	ApplicationID     string
	ApplicationSecret string
	AccessToken       string
	ExpiresAt         int
}

var (
	version       = "0.1.0"
	revision      = "latest"
	tokenFilePath string
)

const tokenFileName = ".mediumctl"

func showPostedArticleInfo(p *medium.PostedArticle) {
	fmt.Println("Your article was successfully posted.\n")
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

func parseArticle(filename string) (article medium.Article, err error) {
	format := "markdown"
	if strings.HasSuffix(filename, "html") || strings.HasSuffix(filename, "htm") {
		format = "html"
	}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	lines := strings.Split(string(b), "\n")
	if len(lines) < 1 {
		err = fmt.Errorf("%s is empty", filename)
		return
	}
	var (
		title        string
		tags         []string
		content      string
		license      string
		status       string
		canonicalURL string
		notify       bool
	)

	title = "untitled"
	status = "public"

	for i, line := range lines[1:] {
		if strings.HasPrefix(line, "---") {
			content = strings.Join(lines[i+1:], "\n")
			break
		}
		if strings.HasPrefix(line, "title: ") {
			title = line[len("title: "):]
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
		if strings.HasPrefix(line, "tags: ") {
			tags = strings.Split(line[len("tags: "):], " ")
		}
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

func saveToken(clientID, clientSecret string, token *medium.Token) (err error) {
	b, err := json.Marshal(Token{
		ApplicationID:     clientID,
		ApplicationSecret: clientSecret,
		AccessToken:       token.AccessToken,
		ExpiresAt:         token.ExpiresAt,
	})
	if err != nil {
		return
	}
	err = ioutil.WriteFile(tokenFilePath, b, 0644)
	return
}

func loadToken() (*Token, error) {
	b, err := ioutil.ReadFile(tokenFilePath)
	if err != nil {
		return nil, fmt.Errorf("API token is not set. Please run 'auth' at first.")
	}
	var token Token
	err = json.Unmarshal(b, &token)
	return &token, err
}

func main() {
	err := Run(os.Args)

	if err != nil {
		log.New(os.Stderr, "error: ", 0).Fatal(err)
		os.Exit(1)
	}

	os.Exit(0)
}
func Run(args []string) (err error) {
	if len(args) < 2 {
		return Help(args)
	}
	u, err := user.Current()
	if err != nil {
		return
	}
	tokenFilePath = filepath.Join(u.HomeDir, tokenFileName)
	switch args[1] {
	case "auth":
		err = Auth(args)
	case "info":
		err = Info(args)
	case "p":
		err = Post(args, false)
	case "publication":
	case "u":
		err = Post(args, true)
	case "user":
		err = Post(args, true)
	case "version":
		err = Version(args)
	case "help":
		err = Help(args)
	default:
		fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s subcommand.\n", args[0], args[1], args[0])
		err = Help(args)
	}
	return
}

func Auth(args []string) (err error) {
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

func Info(args []string) (err error) {
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
	fmt.Println("You have publication(s) below:")
	for _, p := range ps {
		fmt.Printf("Name: %s\n", p.Name)
		fmt.Printf("Description: %s\n", p.Description)
		fmt.Printf("URL: %s\n", p.URL)
	}
	return
}

func Post(args []string, userFlag bool) (err error) {
	var (
		debugFlag bool
	)

	f := flag.NewFlagSet(fmt.Sprintf("%s %s", args[0], args[1]), flag.ExitOnError)
	f.BoolVar(&debugFlag, "debug", false, "Enable debug output.")
	f.Parse(args[2:])

	article, err := parseArticle(f.Args()[0])
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
	p, err := ps[0].Post(article)
	if err != nil {
		return
	}
	showPostedArticleInfo(p)
	return
}

func Version(args []string) (err error) {
	fmt.Printf("%s-%s\n", version, revision)
	return
}

func Help(args []string) (err error) {
	fmt.Println(`usage: mediumctl [command] [options]

Command:
  auth
    Setting up API token for Medium with OAuth.
  info
    Show the information about current user and its publications.
  user, u
    Post HTML or Markdown file to current user profile.
  publication, p
    Post HTML or Markdown file to current user's publication.
  version
  Show version and revision information.
  help
    Show this message.

For more information, please see https://github.com/moutend/mediumctl.`)
	return
}
