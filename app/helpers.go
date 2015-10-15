package main

import (
	"crypto/md5"
	"fmt"

	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/sessions"
)

// PageData represents data for passing in the templates
type PageData struct {
	User          *User
	GravatarImage string
	AppUrl        string
	SuccessFlash  []interface{}
	ErrorFlash    []interface{}
	Errors        *binding.Errors
	Config        *BareCVConfig
}

// NewPageData is the constructor for PageData struct
func NewPageData(tokens oauth2.Tokens, session sessions.Session) *PageData {
	pd := &PageData{}

	pd.User = CurrentUser(tokens)

	pd.GravatarImage = fmt.Sprintf(
		"https://secure.gravatar.com/avatar/%x?s=50&d=https%%3A%%2F%%2Fwww.synkee.com%%2Fapp%%2Fstatic%%2Fdist%%2Fimg%%2Fuser.png",
		md5.Sum([]byte(pd.User.Email)))

	pd.AppUrl = config.AppUrl
	pd.Errors = &binding.Errors{}

	pd.SuccessFlash = session.Flashes("success")
	pd.ErrorFlash = session.Flashes("error")

	pd.Config = config

	return pd
}

// UserCV struct represents user's CV object
type UserCV struct {
	User *User
}

// GeneratePDFHelper generates PDF using rubber
func GeneratePDFHelper(user *User) {

}
