package main

import (
	"github.com/martini-contrib/oauth2"
	_ "log"
)

// CurrentUser returns current signed in user object
func CurrentUser(tokens oauth2.Tokens) *User {
	profile, _ := GoogleGetProfile(tokens.Access())
	email := profile.Email
	name := profile.Name
	// log.Printf("[CurrentUser] User's email: %s", email)

	user := User{Email: email}
	db.Where(&user).First(&user)

	// log.Printf("[CurrentUser] User: %s", user)

	if user.ID == 0 {
		db.Unscoped().Where(&user).First(&user)
	}

	if len(user.Name) == 0 && len(name) > 0 {
		user.Name = name
		db.Save(&user)
	}

	if user.DeletedAt != nil {
		// log.Printf("[CurrentUser] User.DeletedAt: %s", user.DeletedAt)
		user.DeletedAt = nil
		db.Unscoped().Save(&user)
	}

	// log.Printf("[CurrentUser] User: %s", user)
	// log.Printf("[CurrentUser] Token: %s", tokens.Access())

	return &user
}
