package dto

type UserCreateDto struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required"`
}
