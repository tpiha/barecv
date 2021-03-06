package main

import (
	"os"
	"path/filepath"

	"github.com/martini-contrib/oauth2"
)

// CurrentUser returns current signed in user object
func CurrentUser(tokens oauth2.Tokens) *User {
	profile, _ := GoogleGetProfile(tokens.Access())
	email := profile.Email
	name := profile.Name

	user := User{Email: email}
	db.Where(&user).First(&user)

	if user.ID == 0 {
		db.Unscoped().Where(&user).First(&user)
	}

	if len(user.Name) == 0 && len(name) > 0 {
		user.Name = name
		db.Save(&user)

		setting := &Setting{UserID: int(user.ID)}
		db.Create(setting)
	}

	if user.DeletedAt != nil {
		user.DeletedAt = nil
		db.Unscoped().Save(&user)

		if len(user.Username) > 0 {
			path := filepath.Join(HomeDir(), "app/public/files/"+user.Username)
			os.Mkdir(path, 0755)
		}
	}

	return &user
}
