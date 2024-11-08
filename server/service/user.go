package service

import (
	"errors"
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/repository"
	"watt-flow/util"
)

type IUserService interface {
	FindById(id uint64) (*model.User, error)
	Create(user *model.User) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Login(loginCredentials dto.LoginDto) (string, error)
	Register(registrationDto *dto.RegistrationDto) (*model.User, error)
}
type UserService struct {
	repository  *repository.UserRepository
	authService *AuthService
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

func (service *UserService) Login(loginCredentials dto.LoginDto) (string, error) {
	user, err := service.repository.FindByEmail(loginCredentials.Username)
	if err != nil {
		return "", err
	}
	if !util.ComparePasswords(user.Password, loginCredentials.Password) {
		return "", errors.New("Invalid credentials")
	} else {
		token := service.authService.CreateToken(user)
		return token, nil
	}
}

func (service *UserService) Register(registrationDto *dto.RegistrationDto) (*model.User, error) {
	user := model.User{}
	user.Username = registrationDto.Username
	user.Email = registrationDto.Email
	user.Password = util.HashPassword(registrationDto.Password)
	//user.ProfileImage = registrationDto.ProfileImage
	user.Role = 0
	user.Status = 1
	createdUser, err := service.repository.Create(&user)
	if err != nil {
		return nil, err
	}
	return &createdUser, nil
}

func NewUserService(repository *repository.UserRepository, authService *AuthService) *UserService {
	return &UserService{
		repository:  repository,
		authService: authService,
	}
}
