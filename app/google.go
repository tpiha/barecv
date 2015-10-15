package main

import (
	"github.com/martini-contrib/sessions"
	"golang.org/x/oauth2"
	"log"

	martinioauth2 "github.com/martini-contrib/oauth2"
	goauth2 "google.golang.org/api/oauth2/v1"
)

// InitGoogle initializes Google endpoint for oauth2
func InitGoogle() {
	m.Use(sessions.Sessions("tinycv", sessions.NewCookieStore([]byte("Eogoofie7aV0ouyae3ieBeefaij9mohj"))))

	m.Use(martinioauth2.Google(
		&oauth2.Config{
			ClientID:     "716025935670-v79abqso61n9loh6b6n9btktj73cb5qr.apps.googleusercontent.com",
			ClientSecret: "7BiYjPxvKlkiQFJzVWqYNVqo",
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			RedirectURL:  config.AppUrl + "/oauth2callback",
		},
	))
}

// GoogleGetProfile returns Google profile object (oauth2.Userinfoplus)
func GoogleGetProfile(tokenStr string) (*goauth2.Userinfoplus, error) {
	token := &oauth2.Token{AccessToken: tokenStr}
	conf := &oauth2.Config{}
	oauthClient := conf.Client(oauth2.NoContext, token)

	client, err := goauth2.New(oauthClient)
	if err != nil {
		log.Printf("[GoogleGetProfile] error: %s", err)
		return &goauth2.Userinfoplus{}, err
	}

	user, err := client.Userinfo.Get().Do()
	if err != nil {
		log.Printf("[GoogleGetProfile] error: %s", err)
		return &goauth2.Userinfoplus{}, err
	}

	return user, nil
}
