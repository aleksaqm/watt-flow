package service

import (
	"watt-flow/model"
	"watt-flow/repository"
)

type IUserService interface {
	FindById(id string) (*model.User, error)
}
type UserService struct {
	repository repository.UserRepository
}

func (service *UserService) FindById(id string) (*model.User, error) {
	user := model.User{
		FirstName: "danilo",
		LastName:  "cvijetic",
		Id:        312,
	}

	return &user, nil
}
func NewUserService(repository repository.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}
