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
