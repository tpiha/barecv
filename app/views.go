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

// CV renders dashboard page
func CV(r render.Render, tokens oauth2.Tokens, session sessions.Session) {
	pd := NewPageData(tokens, session)
	r.HTML(200, "cv", pd)
}

// CVSave saves dashboard page
func CVSave(r render.Render, tokens oauth2.Tokens, session sessions.Session, profile ProfileForm, err binding.Errors) {
	pd := NewPageData(tokens, session)

	log.Printf("[CVSave] profile: %s", profile)

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
		r.Redirect(config.AppUrl+"/cv", 302)
	} else {
		pd.Errors = &err
		log.Printf("[CVSave] errors: %s", err[0].FieldNames)
		r.HTML(200, "cv", pd)
	}
}

// CVSaveSocial saves dashboard page
func CVSaveSocial(r render.Render, tokens oauth2.Tokens, session sessions.Session, social SocialNetworksForm, err binding.Errors) {
	pd := NewPageData(tokens, session)

	log.Printf("[CVSaveSocial] social: %s", social)

	user := pd.User
	user.LinkedIn = social.LinkedIn
	user.Facebook = social.Facebook
	user.Twitter = social.Twitter
	user.GitHub = social.GitHub
	user.Instagram = social.Instagram
	db.Save(user)
	session.AddFlash("You have successfully updated your BareCV profile.", "success")
	r.Redirect(config.AppUrl+"/cv", 302)
}

// Account renders user's account page
func Account(r render.Render, tokens oauth2.Tokens, session sessions.Session) {
	pd := NewPageData(tokens, session)
	r.HTML(200, "account", pd)
}

// AccountDelete deletes user's SA account
func AccountDelete(r render.Render, tokens oauth2.Tokens, session sessions.Session) {
	pd := NewPageData(tokens, session)
	db.Delete(pd.User)
	r.Redirect(config.AppUrl, 302)
}
