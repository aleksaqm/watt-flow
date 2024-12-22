package dto

type ClerkRegisterDto struct {
	Username     string `json:"username" binding:"required,min=3,max=20"`
	Jmbg         string `json:"jmbg" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	ProfileImage string `json:"profile_image"`
}
