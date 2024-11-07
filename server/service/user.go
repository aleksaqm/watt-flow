package service

import (
	"watt-flow/model"
	"watt-flow/repository"
	"watt-flow/util"
)

type IUserService interface {
	FindById(id uint64) (*model.User, error)
	Create(user *model.User) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
}
type UserService struct {
	repository *repository.UserRepository
}

func (service *UserService) FindById(id uint64) (*model.User, error) {
	user, _ := service.repository.FindById(id)
	return user, nil
}

func (service *UserService) FindByEmail(email string) (*model.User, error) {
	user, err := service.repository.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *UserService) Create(user *model.User) (*model.User, error) {
	user.Password = util.HashPassword(user.Password)
	createdUser, err := service.repository.Create(user)
	if err != nil {
		return nil, err
	}
	return &createdUser, nil
}

//func (service *UserService) Login(loginDto dto.LoginDto) (string, error) {
//	user, err := service.repository.FindByEmail(loginDto.Username)
//	if err != nil {
//		return "", err
//	}
//	hashedPassword := util.HashPassword(loginDto.Password)
//	if user.Password != hashedPassword {
//		return "", errors.New("Invalid credentials")
//	}
//
//}

func NewUserService(repository *repository.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}
