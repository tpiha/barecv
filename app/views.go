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
	pd.Sections = pd.User.GetSections()
	r.HTML(200, "sections", pd)
}

// SectionsNew renders page for adding new CV section
func SectionsNew(r render.Render, tokens oauth2.Tokens, session sessions.Session, params martini.Params) {
	pd := NewPageData(tokens, session)
	sectionType, _ := strconv.Atoi(params["type"])
	pd.SectionType = sectionType
	pd.Section = &Section{}
	r.HTML(200, "sections-new", pd)
}

// SectionsNewPost saves section to the database
func SectionsNewPost(r render.Render, tokens oauth2.Tokens, session sessions.Session, params martini.Params, req *http.Request) {
	pd := NewPageData(tokens, session)

	pd.SectionType, _ = strconv.Atoi(params["type"])
	pd.Section = &Section{}

	var title, subtitle, left, right string

	errors := &binding.Errors{}
	req.ParseForm()

	if pd.SectionType == TypeTitle {
		title = req.Form.Get("title")
		pd.Section.Title = title
		log.Printf("[SectionsNewPost] title: %s", title)
		if len(title) == 0 {
			errors.Add([]string{"title"}, "RequiredError", "This field is required.")
		}
	} else if pd.SectionType == TypeSubtitle {
		subtitle = req.Form.Get("subtitle")
		pd.Section.Subtitle = subtitle
		if len(subtitle) == 0 {
			errors.Add([]string{"subtitle"}, "RequiredError", "This field is required.")
		}
	} else if pd.SectionType == TypeParagraph {
		left = req.Form.Get("left")
		right = req.Form.Get("right")
		pd.Section.Left = left
		pd.Section.Right = right
		if len(left) == 0 {
			errors.Add([]string{"left"}, "RequiredError", "This field is required.")
		}
		if len(right) == 0 {
			errors.Add([]string{"right"}, "RequiredError", "This field is required.")
		}
	}

	if errors.Len() == 0 {
		section := &Section{Type: pd.SectionType, User: *pd.User, OrderID: pd.User.GetLastSectionOrderId()}
		if pd.SectionType == TypeTitle {
			section.Title = title
		} else if pd.SectionType == TypeSubtitle {
			section.Subtitle = subtitle
		} else if pd.SectionType == TypeParagraph {
			section.Left = left
			section.Right = right
		}
		db.Save(section)
		session.AddFlash("You have successfully added a new section to your CV.", "success")
		r.Redirect(config.AppUrl+"/sections", 302)
	} else {
		pd.Errors = errors
		log.Printf("[Save] errors: %s", errors)
		r.HTML(200, "sections-new", pd)
	}
}

// SectionsDelete deletes a section from user's CV
func SectionsDelete(r render.Render, tokens oauth2.Tokens, session sessions.Session, params martini.Params) {
	sectionID, _ := strconv.Atoi(params["section_id"])
	section := &Section{}

	db.Delete(section, sectionID)

	session.AddFlash("You have successfully deleted a section from your CV.", "success")
	r.Redirect(config.AppUrl+"/sections", 302)
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
