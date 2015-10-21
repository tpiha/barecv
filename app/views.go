package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

// Home renders home page
func Home(r render.Render, req *http.Request, tokens oauth2.Tokens, session sessions.Session, params martini.Params) {
	if "http://"+req.Host == config.AppUrl {
		o := render.HTMLOptions{Layout: ""}
		r.HTML(200, "home", nil, o)
	} else {
		ShowPublic(r, req)
	}
}

// Show renders user's CV
func Show(r render.Render, req *http.Request, username string) {
	o := render.HTMLOptions{Layout: ""}
	cv := &UserCV{}

	cv.User = &User{}
	cv.User.Username = username
	db.Where(cv.User).First(cv.User)
	cv.Settings = &Setting{}
	cv.Settings.UserID = int(cv.User.ID)
	db.Where(cv.Settings).First(cv.Settings)

	visit := &Visit{User: *cv.User}
	db.Create(visit)
	log.Printf("[Show] font: %s", cv.Settings.Font)
	r.HTML(200, "cv-templates/default", cv, o)
}

// ShowPublic renders public user's CV
func ShowPublic(r render.Render, req *http.Request) {
	username := strings.Replace(req.Host, config.CVBase, "", -1)

	user := &User{}
	user.Username = username
	db.Where(user).First(user)
	settings := &Setting{}
	settings.UserID = int(user.ID)
	db.Where(settings).First(settings)

	if settings.PrivacyLevel == PrivacyPublic {
		Show(r, req, username)
	} else {
		r.Redirect(config.AppUrl, 302)
	}
}

// ShowPrivate renders private user's CV
func ShowPrivate(r render.Render, req *http.Request, tokens oauth2.Tokens, session sessions.Session, params martini.Params) {
	pd := NewPageData(tokens, session)
	username := params["username"]
	if pd.User.Username == username {
		Show(r, req, username)
	} else {
		session.AddFlash("This is not your CV.", "error")
		r.Redirect(config.AppUrl+"/dashboard", 302)
	}
}

// ShowPassword renders user's CV with password
func ShowPassword() {

}

// ShowHash renders user's CV with hash
func ShowHash(r render.Render, req *http.Request, tokens oauth2.Tokens, session sessions.Session, params martini.Params) {
	username := strings.Replace(req.Host, config.CVBase, "", -1)

	user := &User{}
	user.Username = username
	db.Where(user).First(user)
	settings := &Setting{}
	settings.UserID = int(user.ID)
	db.Where(settings).First(settings)

	if settings.Hash == params["hash"] {
		Show(r, req, username)
	} else {
		r.Redirect(config.AppUrl, 302)
	}
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

// SectionsUpdate renders page for updating CV section
func SectionsUpdate(r render.Render, tokens oauth2.Tokens, session sessions.Session, params martini.Params) {
	pd := NewPageData(tokens, session)
	sectionType, _ := strconv.Atoi(params["type"])
	sectionID, _ := strconv.Atoi(params["section_id"])
	pd.SectionType = sectionType
	pd.Section = &Section{}
	db.First(pd.Section, sectionID)
	r.HTML(200, "sections-update", pd)
}

// SectionsPost saves new section to the database
func SectionsPost(r render.Render, tokens oauth2.Tokens, session sessions.Session, params martini.Params, req *http.Request) {
	pd := NewPageData(tokens, session)
	req.ParseForm()

	pd.SectionType, _ = strconv.Atoi(params["type"])
	pd.Section = &Section{}
	action := req.Form.Get("action")

	log.Printf("[SectionPost] action: %s", action)

	var title, subtitle, left, right string

	errors := &binding.Errors{}
	req.ParseForm()

	if pd.SectionType == TypeTitle {
		title = req.Form.Get("title")
		pd.Section.Title = title
		log.Printf("[SectionsPost] title: %s", title)
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
		section := &Section{}
		if action == "update" {
			sectionID, _ := strconv.Atoi(params["section_id"])
			db.First(section, sectionID)
		} else {
			section.Type = pd.SectionType
			section.User = *pd.User
			section.OrderID = pd.User.GetLastSectionOrderId()
		}
		if pd.SectionType == TypeTitle {
			section.Title = title
		} else if pd.SectionType == TypeSubtitle {
			section.Subtitle = subtitle
		} else if pd.SectionType == TypeParagraph {
			section.Left = left
			section.Right = right
		}
		db.Save(section)
		if action == "update" {
			session.AddFlash("You have successfully updated a section of your CV.", "success")
		} else {
			session.AddFlash("You have successfully added a new section to your CV.", "success")
		}
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
	user.LinkedIn = strings.Replace(strings.Replace(social.LinkedIn, "https://", "", -1), "http://", "", -1)
	user.Facebook = strings.Replace(strings.Replace(social.Facebook, "https://", "", -1), "http://", "", -1)
	user.Twitter = strings.Replace(strings.Replace(social.Twitter, "https://", "", -1), "http://", "", -1)
	user.GitHub = strings.Replace(strings.Replace(social.GitHub, "https://", "", -1), "http://", "", -1)
	user.Instagram = strings.Replace(strings.Replace(social.Instagram, "https://", "", -1), "http://", "", -1)
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
	os.RemoveAll(filepath.Join(HomeDir(), "app/public/files/"+pd.User.Username))
	r.Redirect(config.AppUrl, 302)
}

// AccountRedirect redirects user to the account page with a flash message
func AccountRedirect(r render.Render, tokens oauth2.Tokens, session sessions.Session) {
	session.AddFlash("First you must choose a BareCV username.", "error")
	r.Redirect(config.AppUrl+"/account", 302)
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

// Reorder sections on the CV
func Reorder(r render.Render, tokens oauth2.Tokens, session sessions.Session, req *http.Request) {
	req.ParseForm()
	order := req.Form["order[]"]

	for i, sectionId := range order {
		section := &Section{}
		sId, _ := strconv.Atoi(sectionId)
		section.ID = uint(sId)
		db.First(section, &section.ID)
		log.Printf("[Reorder] section: %s", section)
		section.OrderID = i
		db.Save(section)
	}

	r.JSON(200, map[string]interface{}{"success": true})
}

// Settings renders user's settings page
func Settings(r render.Render, tokens oauth2.Tokens, session sessions.Session) {
	pd := NewPageData(tokens, session)
	pd.Fonts = GetFonts()
	r.HTML(200, "settings", pd)
}

// SettingsSave saves user's settings
func SettingsSave(r render.Render, tokens oauth2.Tokens, session sessions.Session, settings SettingsForm, err binding.Errors, req *http.Request) {
	pd := NewPageData(tokens, session)

	log.Printf("[SettingsSave] settings: %s", settings)
	log.Printf("[SettingsSave] settings.SearchIndexingEnabled: %s", settings.SearchIndexingEnabled)

	userSettings := pd.Settings
	userSettings.Color = settings.Color
	userSettings.Font = settings.Font
	userSettings.GoogleAnalytics = settings.GoogleAnalytics
	userSettings.SearchIndexingEnabled = settings.SearchIndexingEnabled == "on"
	userSettings.PrivacyLevel = settings.PrivacyLevel
	if userSettings.PrivacyLevel == PrivacyHash {
		userSettings.Hash = RandSeq(32)
	}
	db.Save(userSettings)

	session.AddFlash("You have successfully saved your settings.", "success")
	r.Redirect(config.AppUrl+"/settings", 302)
}
