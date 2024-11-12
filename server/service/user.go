package service

import (
	"errors"
	"fmt"
	"strings"
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/repository"
	"watt-flow/util"
)

type IUserService interface {
	FindById(id uint64) (*model.User, error)
	Create(user *dto.UserCreateDto) (*dto.UserDto, error)
	FindByEmail(email string) (*model.User, error)
	Login(loginCredentials dto.LoginDto) (string, error)
	Register(registrationDto *dto.RegistrationDto) (*dto.UserDto, error)
	ActivateAccount(token string) error
	CreateSuperAdmin() (string, error)
	ChangeAdminPassword(passwordDto dto.NewPasswordDto) error
	IsAdminActive() bool
	FindAllByRole(role string) (*[]dto.UserDto, error)
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

func (service *UserService) Create(userDto *dto.UserCreateDto) (*dto.UserDto, error) {
	userDto.Password = util.HashPassword(userDto.Password)
	role, err := model.ParseRole(userDto.Role)
	if err != nil {
		return nil, err
	}
	user := model.User{
		Username: userDto.Username,
		Password: userDto.Password,
		Email:    userDto.Email,
		Role:     role,
		Status:   model.Active,
	}
	createdUser, err := service.repository.Create(&user)
	if err != nil {
		return nil, err
	}
	userReturn := dto.UserDto{
		Id:       createdUser.Id,
		Username: createdUser.Username,
		Email:    createdUser.Email,
		Role:     createdUser.Role.RoleToString(),
	}
	return &userReturn, nil
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

func (service *UserService) Register(registrationDto *dto.RegistrationDto) (*dto.UserDto, error) {
	user := model.User{}
	user.Username = registrationDto.Username
	user.Email = registrationDto.Email
	user.Password = util.HashPassword(registrationDto.Password)
	user.Role = model.Regular
	user.Status = model.Inactive
	if registrationDto.ProfileImage != "" {
		base64String := registrationDto.ProfileImage
		if strings.HasPrefix(base64String, "data:image/") {
			base64String = strings.SplitN(base64String, ",", 2)[1]
		}
		filePath, err := util.SaveFile(user.Username, base64String, "jpg", "profile_pictures")
		if err != nil {
			return nil, err
		}
		user.ProfileImage = filePath
	}
	err := service.SendActivationEmail(&user)
	if err != nil {
		return nil, err
	}
	createdUser, err := service.repository.Create(&user)
	if err != nil {
		return nil, err
	}
	userDto := dto.UserDto{
		Id:       createdUser.Id,
		Username: createdUser.Username,
		Email:    createdUser.Email,
		Role:     createdUser.Role.RoleToString(),
	}
	return &userDto, nil
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
	//emailBody := fmt.Sprintf("<html><body><p>Click <a href='%s'>here</a> to activate your account.</p></body></html>", activationLink)
	emailBody := util.GenerateActivationEmailBody(activationLink)
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
		Status:   model.Inactive,
		Role:     model.SuperAdmin,
	}
	_, err = service.repository.Create(&superAdmin)
	if err != nil {
		return "", errors.New("error creating superAdmin")
	}
	return password, nil
}

func (service *UserService) ChangeAdminPassword(passwordDto dto.NewPasswordDto) error {
	admin, err := service.repository.FindByUsername("admin")
	if admin.Status == model.Active {
		return errors.New("admin is already active")
	}
	if err != nil {
		return err
	}
	if !util.ComparePasswords(admin.Password, passwordDto.OldPassword) {
		return errors.New("invalid credentials")
	}
	admin.Password = util.HashPassword(passwordDto.NewPassword)
	admin.Status = model.Active
	_, err = service.repository.Update(admin)
	if err != nil {
		return err
	}
	return nil
}

func (service *UserService) IsAdminActive() bool {
	admin, err := service.repository.FindByUsername("admin")
	if err != nil {
		return false
	}
	return admin.Status == model.Active
}

func (service *UserService) FindAllByRole(roleStr string) (*[]dto.UserDto, error) {
	role, err := model.ParseRole(roleStr)
	if err != nil {
		return nil, err
	}
	users, err := service.repository.FindAllByRole(role)
	if err != nil {
		return nil, err
	}
	usersDto := make([]dto.UserDto, len(*users))
	for i, user := range *users {
		usersDto[i] = dto.UserDto{
			Id:       user.Id,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role.RoleToString(),
		}
	}
	return &usersDto, nil
}

func NewUserService(repository *repository.UserRepository, authService *AuthService) *UserService {
	return &UserService{
		repository:  repository,
		authService: authService,
	}
}
