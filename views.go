package main

import (
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

// Home renders home page
func Home(r render.Render) {
	o := render.HTMLOptions{Layout: ""}
	r.HTML(200, "home", nil, o)
}

// Dashboard renders dashboard page
func Dashboard(r render.Render, tokens oauth2.Tokens, session sessions.Session) {
	pd := NewPageData(tokens, session)
	r.HTML(200, "dashboard", pd)
}
