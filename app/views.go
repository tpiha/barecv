package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

// Home renders home page
func Home(r render.Render, req *http.Request) {
	if "http://"+req.Host == config.AppUrl {
		o := render.HTMLOptions{Layout: ""}
		r.HTML(200, "home", nil, o)
	} else {
		Show(r, req)
	}
}

// Show renders user's CV
func Show(r render.Render, req *http.Request) {
	o := render.HTMLOptions{Layout: ""}
	cv := &UserCV{}
	username := strings.Replace(req.Host, config.CVBase, "", -1)
	log.Printf("[Show] URL: %s", req.Host)
	log.Printf("[Show] username: %s", username)
	cv.User = &User{}
	cv.User.Username = username
	db.Where(cv.User).First(cv.User)

	visit := &Visit{User: *cv.User}
	db.Create(visit)

	r.HTML(200, "cv-templates/default", cv, o)
}

// Dashboard renders dashboard page
func Dashboard(r render.Render, tokens oauth2.Tokens, session sessions.Session) {
	pd := NewPageData(tokens, session)
	pd.ChartData = GetChartData(pd.User)
	r.HTML(200, "dashboard", pd)
}

// CV renders CV page
func CV(r render.Render, tokens oauth2.Tokens, session sessions.Session) {
	pd := NewPageData(tokens, session)
	r.HTML(200, "cv", pd)
}

// Sections renders page for editing CV sections
func Sections(r render.Render, tokens oauth2.Tokens, session sessions.Session) {
	pd := NewPageData(tokens, session)
	r.HTML(200, "sections", pd)
}

// SectionsNew renders page for adding new CV section
func SectionsNew(r render.Render, tokens oauth2.Tokens, session sessions.Session, params martini.Params) {
	pd := NewPageData(tokens, session)
	sectionType, _ := strconv.Atoi(params["type"])
	pd.SectionType = sectionType
	r.HTML(200, "sections-new", pd)
}

// Save saves dashboard page
func Save(r render.Render, tokens oauth2.Tokens, session sessions.Session, profile ProfileForm, err binding.Errors) {
	pd := NewPageData(tokens, session)

	log.Printf("[Save] profile: %s", profile)

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
		log.Printf("[Save] errors: %s", err[0].FieldNames)
		r.HTML(200, "cv", pd)
	}
}

// SaveSocial saves dashboard page
func SaveSocial(r render.Render, tokens oauth2.Tokens, session sessions.Session, social SocialNetworksForm, err binding.Errors) {
	pd := NewPageData(tokens, session)

	log.Printf("[SaveSocial] social: %s", social)

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

// AccountSave renders user's account page
func AccountSave(r render.Render, tokens oauth2.Tokens, session sessions.Session, username UsernameForm, err binding.Errors) {
	pd := NewPageData(tokens, session)

	if len(username.Username) > 0 {
		existingUser := &User{Username: username.Username}
		db.Where(existingUser).First(existingUser)

		if existingUser.ID > 0 {
			err.Add([]string{"username"}, "RequiredError", "This username is already taken.")
		}

		log.Printf("[AccountSave] existing user: %s", existingUser)
	}

	if err.Len() == 0 {
		user := pd.User
		user.Username = username.Username
		db.Save(user)
		session.AddFlash("You have successfully updated your BareCV username / domain.", "success")
		r.Redirect(config.AppUrl+"/account", 302)
	} else {
		pd.Errors = &err
		log.Printf("[AccountSave] errors: %s", err[0].Classification)
		r.HTML(200, "account", pd)
	}
}

// GeneratePDF generates PDF of user's CV
func GeneratePDF(r render.Render, tokens oauth2.Tokens, session sessions.Session) {
	pd := NewPageData(tokens, session)
	url := "http://" + pd.User.Username + config.CVBase + "/pdf/" + pd.User.Username + ".pdf"
	GeneratePDFHelper(pd.User)
	r.Redirect(url, 302)
}
