package main

import (
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/url"

	"github.com/go-martini/martini"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
)

// Martini object
var m *martini.ClassicMartini

// Config object
var config *BareCVConfig

// Database object
var db *gorm.DB

// Main function
func main() {
	// Initializa config
	config = &BareCVConfig{}
	config.Load("../config.json")

	// Initizalize database
	db = InitDb()

	// Initialize Martini
	m = martini.Classic()

	// Render html templates from templates directory
	m.Use(render.Renderer(render.Options{
		Extensions: []string{".html"},
		Layout:     "layout",
		Funcs: []template.FuncMap{
			{
				"unescaped": func(args ...interface{}) template.HTML {
					log.Printf("[main] unescaped: %s", args[0].(string))
					urlStr, _ := url.QueryUnescape(args[0].(string))
					return template.HTML(urlStr)
				},
			},
		},
	}))

	// Initializes Google auth
	InitGoogle()

	// GET methods
	m.Get("/", Home)
	m.Get("/dashboard", oauth2.LoginRequired, Dashboard)
	m.Get("/cv", binding.Form(ProfileForm{}), oauth2.LoginRequired, CV)
	m.Get("/sections", oauth2.LoginRequired, Sections)
	m.Get("/sections/new/:type", oauth2.LoginRequired, SectionsNew)
	m.Get("/sections/delete/:section_id", oauth2.LoginRequired, SectionsDelete)
	m.Get("/sections/:type/:section_id", oauth2.LoginRequired, SectionsUpdate)
	m.Get("/account", oauth2.LoginRequired, Account)
	m.Get("/account-delete", oauth2.LoginRequired, AccountDelete)
	m.Get("/generate-pdf", oauth2.LoginRequired, GeneratePDF)
	m.Get("/settings", oauth2.LoginRequired, Settings)

	// POST methods
	m.Post("/cv-save", binding.Form(ProfileForm{}), oauth2.LoginRequired, Save)
	m.Post("/cv-save-social", binding.Form(SocialNetworksForm{}), oauth2.LoginRequired, SaveSocial)
	m.Post("/account-save", binding.Form(UsernameForm{}), oauth2.LoginRequired, AccountSave)
	m.Post("/sections/new/:type", oauth2.LoginRequired, SectionsPost)
	m.Post("/sections/:type/:section_id", oauth2.LoginRequired, SectionsPost)
	m.Post("/sections/reorder", oauth2.LoginRequired, Reorder)
	m.Post("/settings", binding.Form(SettingsForm{}), oauth2.LoginRequired, SettingsSave)

	// Run server
	// m.RunOnAddr(":8080")
	m.Run()
}
