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
}
