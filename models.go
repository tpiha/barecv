package main

import (
	"github.com/jinzhu/gorm"
)

// Struct representing user model
type User struct {
	gorm.Model
	Hash           string `sql:"size:255"`
	Name           string `sql:"size:255"`
	Email          string `sql:"size:255"`
	AddressLine1   string `sql:"size:255"`
	AddressLine2   string `sql:"size:255"`
	AddressLine3   string `sql:"size:255"`
	GitHubUsername string `sql:"size:255"`
}
