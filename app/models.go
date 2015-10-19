package main

import (
	"html/template"

	"github.com/jinzhu/gorm"
)

// User represents user database model
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
	Username   string `sql:"size:255"`
}

// GetLastSectionOrderId gets order id for newly created section
func (u *User) GetLastSectionOrderId() int {
	sections := []*Section{}
	db.Where(Section{UserID: int(u.ID)}).Find(&sections)
	return len(sections)
}

// GetSections returns user's sections ordered by ID
func (u *User) GetSections() []*Section {
	sections := []*Section{}
	db.Where(Section{UserID: int(u.ID)}).Order("order_id").Find(&sections)
	return sections
}

// Visit represents CV visit model
type Visit struct {
	gorm.Model
	User   User
	UserID int `sql:"index"`
}

const (
	TypeTitle     = 1
	TypeSubtitle  = 2
	TypeParagraph = 3
)

// Section represents one horizontal section of the CV
type Section struct {
	gorm.Model
	User     User
	UserID   int `sql:"index"`
	Type     int
	OrderID  int
	Title    string `sql:"size:255"`
	Subtitle string `sql:"size:255"`
	Left     string `sql:"size:255"`
	Right    string `sql:"type:text"`
}

// String returns string representation of the section object
func (s *Section) String() string {
	if s.Type == TypeTitle {
		return s.Title
	} else if s.Type == TypeSubtitle {
		return s.Subtitle
	} else if s.Type == TypeParagraph {
		return s.Left
	}
	return ""
}

// TypeString returns section type as a string
func (s *Section) TypeString() string {
	sType := ""

	if s.Type == TypeTitle {
		sType = "title"
	} else if s.Type == TypeSubtitle {
		sType = "subtitle"
	} else if s.Type == TypeParagraph {
		sType = "paragraph"
	}

	return sType
}

// GetRight returns HTML for the right part of the paragraph
func (s *Section) GetRight() template.HTML {
	return template.HTML(s.Right)
}

const (
	PrivacyPublic   = 1
	PrivacyPrivate  = 2
	PrivacyPassword = 3
	PrivacyHash     = 4
)

// Setting represents settings for the users
type Setting struct {
	gorm.Model
	User                  User
	UserID                int    `sql:"index"`
	Color                 string `sql:"size:255;default:'#E80000'"`
	Font                  string `sql:"size:255;default:'Roboto'"`
	PrivacyLevel          int    `sql:"default:1"`
	Template              string `sql:"size:255"`
	GoogleAnalytics       string `sql:"size:255"`
	SearchIndexingEnabled bool   `sql:"is_default:true"`
}
