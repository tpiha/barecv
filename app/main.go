package main

import (
	_ "github.com/lib/pq"

	"github.com/go-martini/martini"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
)

// Martini object
var m *martini.ClassicMartini

// Config object
var config *SaConfig

// Database object
var db *gorm.DB

// Main function
func main() {
	// Initializa config
	config = &SaConfig{}
	config.Load("../config.json")

	// Initizalize database
	db = InitDb()

	// Initialize Martini
	m = martini.Classic()

	// Render html templates from templates directory
	m.Use(render.Renderer(render.Options{
		Extensions: []string{".html"},
		Layout:     "layout",
	}))

	// Initializes Google auth
	InitGoogle()

	// GET methods
	m.Get("/", Home)
	m.Get("/dashboard", binding.Form(ProfileForm{}), oauth2.LoginRequired, Dashboard)
	m.Get("/cv", binding.Form(ProfileForm{}), oauth2.LoginRequired, CV)

	// POST methods
	m.Post("/cv-save", binding.Form(ProfileForm{}), oauth2.LoginRequired, CVSave)

	// Run server
	m.Run()
}
