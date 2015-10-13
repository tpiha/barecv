package main

import (
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
	"log"
)

// Renders home page
func Home(r render.Render) {
	o := render.HTMLOptions{Layout: ""}
	r.HTML(200, "home", nil, o)
}

// Renders dashboard page
func Dashboard(r render.Render, tokens oauth2.Tokens) {
	profile, _ := GitHubGetProfile(tokens.Access())
	log.Printf("[Dashboard] %s", *profile.Name)
	r.HTML(200, "dashboard", nil)
}
