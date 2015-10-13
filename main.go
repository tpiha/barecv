package main

import (
	"github.com/go-martini/martini"
	_ "github.com/lib/pq"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	goauth2 "golang.org/x/oauth2"
)

// Main function
func main() {
	// Initizalize database
	InitDb()

	// Initialize Martini
	m := martini.Classic()

	// Render html templates from templates directory
	m.Use(render.Renderer(render.Options{
		Extensions: []string{".html"},
		Layout:     "layout",
	}))

	m.Use(sessions.Sessions("my_session", sessions.NewCookieStore([]byte("secret123"))))

	m.Use(oauth2.Github(
		&goauth2.Config{
			ClientID:     "96aa9b3f1c9b7ac26c4e",
			ClientSecret: "4c48d5e04a239396c2f4f854b8b500fdfff26de1",
			Scopes:       []string{"user:email"},
			// RedirectURL:  "http://localhost:3000/register-github",
		},
	))

	// GET methods
	m.Get("/", Home)
	m.Get("/dashboard", oauth2.LoginRequired, Dashboard)

	// POST methods

	// Run server
	m.Run()
}
