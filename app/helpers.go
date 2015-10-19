package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

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
	ChartData     []int
	SectionType   int
	Section       *Section
	Sections      []*Section
	Settings      *Setting
	Fonts         []*Font
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

	pd.Settings = &Setting{}
	db.Where(&Setting{UserID: int(pd.User.ID)}).First(pd.Settings)

	log.Printf("[NewPageData] settings: %s", pd.Settings)

	return pd
}

// UserCV struct represents user's CV object
type UserCV struct {
	User     *User
	Settings *Setting
}

// HomeDir returns application home directory
func HomeDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	dir = filepath.Join(dir, "../")
	return dir
}

// GeneratePDFHelper generates PDF using rubber
func GeneratePDFHelper(user *User) {
	dst := filepath.Join(HomeDir(), "temp", user.Username)
	os.MkdirAll(dst, 0755)
	cmd := fmt.Sprintf("cd %s && %s --pdf ../../latex/template.tex", dst, config.RubberBin)
	log.Printf("[GeneratePDFHelper] cmd: %s", cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Printf("[GeneratePDFHelper] error: %s", err)
	}
	log.Printf("[GeneratePDFHelper] out: %s", out)
	src := filepath.Join(dst, "template.pdf")
	newSrc := filepath.Join(HomeDir(), "app/public/pdf", user.Username+".pdf")
	os.Rename(src, newSrc)
	log.Printf("[GeneratePDFHelper] src: %s %s", src, newSrc)
	os.RemoveAll(dst)
}

// GetChartData returns string containing elements for drawing chart on dashboard
func GetChartData(user *User) []int {
	hours := time.Duration(time.Now().Round(time.Hour).Hour()) * time.Hour
	days := time.Duration((time.Now().Day()-1)*24) * time.Hour
	start := time.Now().Round(time.Hour).Add(-hours).Add(-days).AddDate(-1, 1, 0)

	end := start.AddDate(0, 1, 0)

	log.Printf("[GetChartData] start: %s %s", start, end)

	var data []int

	for i := 12; i > 0; i-- {
		var visits []*Visit
		db.Where("created_at BETWEEN ? AND ? AND user_id = ?", start, end, user.ID).Find(&visits)
		data = append(data, len(visits))
		start = start.AddDate(0, 1, 0)
		end = end.AddDate(0, 1, 0)
	}

	return data
}

// type Fonts struct {
// 	Fonts []Font
// }

type Font struct {
	CSSName    string `json:"css-name,omitempty"`
	FontFamily string `json:"font-family,omitempty"`
	FontName   string `json:"font-name,omitempty"`
}

func GetFonts() []*Font {
	content, err := ioutil.ReadFile(filepath.Join(HomeDir(), "fonts.json"))
	if err != nil {
		fmt.Print("Error:", err)
	}

	var fonts []*Font

	err = json.Unmarshal(content, &fonts)
	if err != nil {
		log.Print("Error:", err)
	}

	return fonts
}
