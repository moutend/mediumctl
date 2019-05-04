package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	medium "github.com/moutend/go-medium"
	"github.com/spf13/cobra"
)

var (
	authCommandClientId     string
	authCommandClientSecret string
	authCommandRedirectURL  string
)

var authCommand = &cobra.Command{
	Use:     "auth",
	Short:   "Setup the API token with OAuth",
	Long:    "Setup the API token with OAuth",
	Aliases: []string{"a"},
	RunE: func(c *cobra.Command, args []string) (err error) {
		redirectURL, err := url.Parse(authCommandRedirectURL)
		if err != nil {
			return err
		}

		listener, err := net.Listen("tcp", redirectURL.Hostname()+":"+redirectURL.Port())
		if err != nil {
			return err
		}
		defer listener.Close()

		type response struct {
			Code  string
			Error string
		}

		responseChan := make(chan response)

		go http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte(`<script>window.open("about:blank","_self").close()</script>`))
			w.(http.Flusher).Flush()

			responseChan <- response{
				Code:  req.FormValue("code"),
				Error: req.FormValue("error"),
			}
		}))

		stateBytes := make([]byte, 88)
		_, err = rand.Read(stateBytes)
		if err != nil {
			return err
		}

		code := ""
		state := fmt.Sprintf("%x", stateBytes)
		scope := "basicProfile,listPublications,publishPost"
		uri := fmt.Sprintf("https://medium.com/m/oauth/authorize?client_id=%s&scope=%s&state=%s&response_type=code&redirect_uri=%s", authCommandClientId, scope, state, redirectURL)

		fmt.Println("Please open this URL:", uri)

		select {
		case res := <-responseChan:
			if res.Error != "" {
				return fmt.Errorf("%s", res.Error)
			}

			code = res.Code

			break
		case <-time.After(5 * time.Minute):
			return fmt.Errorf("timeout")
		}

		client = medium.NewClient(authCommandClientId, authCommandClientSecret, "")
		if debug {
			client.SetLogger(log.New(os.Stdout, "Debug: ", 0))
		}

		token, err := client.Token(code, authCommandRedirectURL)
		if err != nil {
			return err
		}
		if err := writeToken(authCommandClientId, authCommandClientSecret, token); err != nil {
			return err
		}

		fmt.Println("Done")

		return nil
	},
}

func init() {
	authCommand.Flags().StringVarP(&authCommandClientId, "id", "i", "", "specify client id")
	authCommand.Flags().StringVarP(&authCommandClientSecret, "secret", "s", "", "specify client secret")
	authCommand.Flags().StringVarP(&authCommandRedirectURL, "redirect", "r", "", "specify redirect url (e.g. 127.0.0.1)")

	rootCommand.AddCommand(authCommand)
}
