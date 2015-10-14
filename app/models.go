package main

import (
	"github.com/jinzhu/gorm"
)

// Struct representing user model
type User struct {
	gorm.Model
	Hash       string `sql:"size:255"`
	Name       string `sql:"size:255"`
	Profession string `sql:"size:255"`
	Email      string `sql:"size:255"`
	Phone      string `sql:"size:255"`
	Website    string `sql:"size:255"`
	Address    string `sql:"size:1023"`
	LinkedIn   string `sql:"size:255"`
	Facebook   string `sql:"size:255"`
	Twitter    string `sql:"size:255"`
	GitHub     string `sql:"size:255"`
	Instagram  string `sql:"size:255"`
}
