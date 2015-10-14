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

	// POST methods
	m.Post("/dashboard-save", binding.Form(ProfileForm{}), oauth2.LoginRequired, DashboardSave)

	// Run server
	m.Run()
}
