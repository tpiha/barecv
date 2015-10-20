package main

// ProfileForm represents form for submitting user's profile to the database
type ProfileForm struct {
	Name       string `form:"name" binding:"required"`
	Profession string `form:"profession" binding:"required"`
	Email      string `form:"email" binding:"required"`
	Phone      string `form:"phone"`
	Website    string `form:"website"`
	Address    string `form:"address"`
}

// SocialNetworksForm represents form for submitting user's profile to the database
type SocialNetworksForm struct {
	LinkedIn  string `form:"linkedin"`
	Facebook  string `form:"facebook"`
	Twitter   string `form:"twitter"`
	GitHub    string `form:"github"`
	Instagram string `form:"instagram"`
}

// UsernameForm represents form for submitting user's profile to the database
type UsernameForm struct {
	Username string `form:"username" binding:"required"`
}

// TitleForm represents form for submitting section tite
type TitleForm struct {
	Title string `form:"title" binding:"required"`
}

// SettingsForm represents form for submitting section tite
type SettingsForm struct {
	Color           string `form:"color" binding:"required"`
	Font            string `form:"font" binding:"required"`
	GoogleAnalytics string `form:"google_analytics"`
}
