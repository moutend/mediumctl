package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	medium "github.com/moutend/go-medium"
)

const (
	TokenFileName = ".mediumctl"
	MediumctlHome = "MEDIUMCTL_HOME"
)

type Token struct {
	ApplicationID     string
	ApplicationSecret string
	AccessToken       string
	ExpiresAt         int
	RefreshToken      string
}

func getTokenPath() (tokenPath string, err error) {
	if tokenPath = os.Getenv(MediumctlHome); tokenPath != "" {
		return tokenPath, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return tokenPath, err
	}

	tokenPath = filepath.Join(wd, TokenFileName)

	return tokenPath, nil
}

func readToken() (token *Token, err error) {
	tokenPath, err := getTokenPath()
	if err != nil {
		return token, err
	}

	file, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return token, fmt.Errorf("API token is not set")
	}
	if err := json.Unmarshal(file, &token); err != nil {
		return token, err
	}

	return token, nil
}

func writeToken(clientID, clientSecret string, token *medium.Token) (err error) {
	file, err := json.Marshal(Token{
		ApplicationID:     clientID,
		ApplicationSecret: clientSecret,
		AccessToken:       token.AccessToken,
		ExpiresAt:         token.ExpiresAt,
		RefreshToken:      token.RefreshToken,
	})
	if err != nil {
		return err
	}

	tokenPath, err := getTokenPath()
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(tokenPath, file, 0644); err != nil {
		return err
	}

	return nil
}
