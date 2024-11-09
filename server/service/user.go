package service

import (
	"errors"
	"fmt"
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
	ActivateAccount(token string) error
	CreateSuperAdmin() (string, error)
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
		user, err = service.repository.FindByUsername(loginCredentials.Username)
		if err != nil || user == nil {
			return "", errors.New("invalid credentials")
		}
	}
	if user.Status == model.Inactive {
		return "", errors.New("user is inactive")
	}
	if !util.ComparePasswords(user.Password, loginCredentials.Password) {
		return "", errors.New("invalid credentials")
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
	user.Role = model.Regular
	user.Status = model.Inactive
	err := service.SendActivationEmail(&user)
	if err != nil {
		return nil, err
	}
	createdUser, err := service.repository.Create(&user)
	if err != nil {
		return nil, err
	}
	return &createdUser, nil
}

func (service *UserService) ActivateAccount(token string) error {
	valid, claims, err := service.authService.Authorize(token)
	if err != nil {
		return err
	}
	if valid {
		if claims != nil {
			email := claims["email"].(string)
			user, err := service.repository.FindByEmail(email)
			if err != nil {
				return err
			}
			user.Status = 0
			_, err2 := service.repository.Update(user)
			if err2 != nil {
				return err2
			}
			return nil
		}
	}
	return errors.New("invalid token")
}

func (service *UserService) SendActivationEmail(user *model.User) error {
	activationToken := service.authService.CreateActivationToken(user)
	activationLink := fmt.Sprintf("http://localhost:5000/activate/%s", activationToken)
	emailBody := fmt.Sprintf("<html><body><p>Click <a href='%s'>here</a> to activate your account.</p></body></html>", activationLink)
	err := util.SendEmail(user.Email, "Activate your account", emailBody)
	return err
}

func (service *UserService) CreateSuperAdmin() (string, error) {
	admin, err := service.repository.FindByUsername("admin")
	if err == nil && admin != nil {
		return "", errors.New("superAdmin already exists")
	}
	password := util.GenerateRandomPassword(16)
	superAdmin := model.User{
		Username: "admin",
		Password: util.HashPassword(password),
		Status:   model.Active,
		Role:     model.SuperAdmin,
	}
	_, err = service.repository.Create(&superAdmin)
	if err != nil {
		return "", errors.New("error creating superAdmin")
	}
	return password, nil
}

func NewUserService(repository *repository.UserRepository, authService *AuthService) *UserService {
	return &UserService{
		repository:  repository,
		authService: authService,
	}
}
