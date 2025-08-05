package service

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/repository"
	"watt-flow/util"

	"gorm.io/gorm"
)

type IUserService interface {
	FindById(id uint64) (*dto.UserDto, error)
	Create(user *dto.UserCreateDto) (*dto.UserDto, error)
	FindByEmail(email string) (*model.User, error)
	Login(loginCredentials dto.LoginDto) (string, error)
	Register(registrationDto *dto.RegistrationDto) (*dto.UserDto, error)
	RegisterClerk(registrationDto *dto.ClerkRegisterDto) (*dto.UserDto, error)
	ActivateAccount(token string) error
	CreateSuperAdmin() (string, error)
	ChangeAdminPassword(passwordDto dto.NewPasswordDto) error
	IsAdminActive() bool
	FindAllByRole(role string) (*[]dto.UserDto, error)
	Query(queryParams *dto.UserQueryParams) ([]dto.UserDto, int64, error)
	Suspend(id uint64) error
	SuspendClerk(id uint64) error
	Unsuspend(id uint64) error
	WithTrx(trxHandle *gorm.DB) IUserService
}

type UserService struct {
	repository     *repository.UserRepository
	meetingService IMeetingService
	emailSender    *util.EmailSender
	authService    *AuthService
}

func (service *UserService) WithTrx(trxHandle *gorm.DB) IUserService {
	return &UserService{
		repository:     service.repository.WithTrx(trxHandle),
		meetingService: service.meetingService.WithTrx(trxHandle),
		emailSender:    service.emailSender,
		authService:    service.authService,
	}
}

func (service *UserService) FindById(id uint64) (*dto.UserDto, error) {
	user, _ := service.repository.FindById(id)
	userReturn := dto.UserDto{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role.RoleToString(),
	}
	return &userReturn, nil
}

func (service *UserService) SuspendClerk(id uint64) error {
	err := service.Suspend(id)
	if err != nil {
		return err
	}
	err = service.meetingService.CancelMeetingsForClerk(id)
	if err != nil {
		return err
	}
	return nil
}

func (service *UserService) Suspend(id uint64) error {
	err := service.repository.ChangeStatus(id, 2)
	if err != nil {
		return err
	}
	return nil
}

func (service *UserService) Unsuspend(id uint64) error {
	err := service.repository.ChangeStatus(id, 0)
	if err != nil {
		return err
	}
	return nil
}

func (service *UserService) FindByEmail(email string) (*model.User, error) {
	user, err := service.repository.FindActiveByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *UserService) Create(userDto *dto.UserCreateDto) (*dto.UserDto, error) {
	existingUser, _ := service.repository.FindByUsernameOrActiveEmail(userDto.Username, userDto.Email)
	if existingUser != nil {
		return nil, errors.New("username or email already taken")
	}
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
	user, err := service.repository.FindForLogin(loginCredentials.Username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if user.Status == model.Inactive {
		return "", errors.New("user is inactive")
	}
	if user.Status == model.Suspended {
		return "", errors.New("user is suspended")
	}
	if !util.ComparePasswords(user.Password, loginCredentials.Password) {
		return "", errors.New("invalid credentials")
	}

	token := service.authService.CreateToken(user)
	return token, nil
}

func (service *UserService) Register(registrationDto *dto.RegistrationDto) (*dto.UserDto, error) {
	existingUser, _ := service.repository.FindByUsername(registrationDto.Username)
	if existingUser != nil {
		return nil, errors.New("username already taken")
	}
	existingUser, _ = service.repository.FindActiveByEmail(registrationDto.Email)
	if existingUser != nil {
		return nil, errors.New("already have account with this email")
	}
	user := model.User{}
	user.Username = registrationDto.Username
	user.Email = registrationDto.Email
	user.FirstName = registrationDto.FirstName
	user.LastName = registrationDto.LastName
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

func (service *UserService) RegisterClerk(registrationDto *dto.ClerkRegisterDto) (*dto.UserDto, error) {
	existingUser, _ := service.repository.FindByUsername(registrationDto.Username)
	if existingUser != nil {
		return nil, errors.New("username already taken")
	}
	existingUser, _ = service.repository.FindActiveByEmail(registrationDto.Email)
	if existingUser != nil {
		return nil, errors.New("already have account with this email")
	}
	user := model.User{}
	user.Username = registrationDto.Username
	user.FirstName = registrationDto.FirstName
	user.LastName = registrationDto.LastName
	user.Email = registrationDto.Email
	user.Password = util.HashPassword(registrationDto.Jmbg)
	user.Role = model.Clerk
	user.Status = model.Active
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
	createdUser, err := service.repository.Create(&user)
	if err != nil {
		return nil, err
	}
	userDto := dto.UserDto{
		Id:        createdUser.Id,
		Username:  createdUser.Username,
		FirstName: createdUser.FirstName,
		LastName:  createdUser.LastName,
		Email:     createdUser.Email,
		Role:      createdUser.Role.RoleToString(),
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
			existingUser, _ := service.repository.FindActiveByEmail(email)
			if existingUser != nil {
				return errors.New("user is already active")
			}
			username := claims["username"].(string)
			user, err := service.repository.FindByEmailAndUsername(email, username)
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
	activationLink := fmt.Sprintf("http://localhost:5000/api/activate/%s", activationToken)
	emailBody := util.GenerateActivationEmailBody(activationLink)
	err := service.emailSender.SendEmail(user.Email, "Activate your account", emailBody)
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

func (service *UserService) Query(queryParams *dto.UserQueryParams) ([]dto.UserDto, int64, error) {
	var users []dto.UserDto
	if queryParams.Search.Id != 0 {
		user, err := service.FindById(queryParams.Search.Id)
		if err != nil {
			return nil, 0, err
		}
		users = make([]dto.UserDto, 0)
		if user != nil {
			users = append(users, *user)

			return users, 1, nil
		}

		return users, 0, nil
	}

	data, count, err := service.repository.Query(queryParams)
	if err != nil {
		log.Printf("Error on query: %v", err)
		return nil, 0, err
	}
	users = make([]dto.UserDto, 0)
	for _, user := range data {
		mapped_user, _ := MapToDto(&user)
		users = append(users, mapped_user)
	}

	return users, count, nil
}

func MapToDto(user *model.User) (dto.UserDto, error) {
	response := dto.UserDto{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role.RoleToString(),
		Username:  user.Username,
		Status:    user.Status.StatusToString(),
	}
	return response, nil
}

func NewUserService(repository *repository.UserRepository, authService *AuthService, emailSender *util.EmailSender, meetingService IMeetingService) *UserService {
	return &UserService{
		repository:     repository,
		authService:    authService,
		emailSender:    emailSender,
		meetingService: meetingService,
	}
}
