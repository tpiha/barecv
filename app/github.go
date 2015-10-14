package main

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
)

func GitHubGetProfile(tokenStr string) (*github.User, error) {
	token := &oauth2.Token{AccessToken: tokenStr}
	conf := &oauth2.Config{}
	oauthClient := conf.Client(oauth2.NoContext, token)

	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get("")
	if err != nil {
		log.Printf("[GitHubGetProfile] error: %s", err)
		return &github.User{}, err
	}

	return user, nil
}
