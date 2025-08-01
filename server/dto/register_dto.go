package dto

type RegistrationDto struct {
	Username     string `json:"username" binding:"required,min=3,max=20"`
	FirstName    string `json:"first_name" binding:"required,min=1,max=50"`
	LastName     string `json:"last_name" binding:"required,min=1,max=50"`
	Password     string `json:"password" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	ProfileImage string `json:"profile_image"`
}
