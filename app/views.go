package main

import (
	"log"

	"github.com/martini-contrib/binding"
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

// DashboardSave saves dashboard page
func DashboardSave(r render.Render, tokens oauth2.Tokens, session sessions.Session, profile ProfileForm, err binding.Errors) {
	pd := NewPageData(tokens, session)

	log.Printf("[DashboardSave] profile: %s", profile)

	if err.Len() == 0 {
		user := pd.User
		user.Name = profile.Name
		user.Profession = profile.Profession
		user.Email = profile.Email
		user.Phone = profile.Phone
		user.Website = profile.Website
		user.Address = profile.Address
		db.Save(user)
		session.AddFlash("You have successfully updated your BareCV profile.", "success")
		r.Redirect(config.AppUrl+"/dashboard", 302)
	} else {
		pd.Errors = &err
		log.Printf("[DashboardSave] errors: %s", err[0].FieldNames)
		r.HTML(200, "dashboard", pd)
	}
}
