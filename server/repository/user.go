package repository

import (
	"watt-flow/db"
	"watt-flow/model"
	"watt-flow/util"
)

type UserRepository struct {
	database db.Database
	logger   util.Logger
}

func NewUserRepository(db db.Database, logger util.Logger) *UserRepository {
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		logger.Error("Error migrating user", err)
	}
	return &UserRepository{
		database: db,
		logger:   logger,
	}
}

func (repository *UserRepository) Create(user *model.User) (model.User, error) {
	result := repository.database.Create(user)
	if result.Error != nil {
		repository.logger.Error("Error creating user", result.Error)
		return *user, result.Error
	}
	return *user, nil
}

func (repository *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	result := repository.database.Where("email = ?", email).First(&user)
	if result.Error != nil {
		repository.logger.Error("Error finding user by email", result.Error)
		return nil, result.Error
	}
	return &user, nil
}

func (repository *UserRepository) FindById(id uint64) (*model.User, error) {
	var user model.User
	if err := repository.database.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repository *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	result := repository.database.Where("username = ?", username).First(&user)
	if result.Error != nil {
		repository.logger.Error("Error finding user by username", result.Error)
		return nil, result.Error
	}
	return &user, nil
}

func (repository *UserRepository) Update(user *model.User) (model.User, error) {
	result := repository.database.Save(user)
	if result.Error != nil {
		repository.logger.Error("Error updating user", result.Error)
		return *user, result.Error
	}
	return *user, nil
}

func (repository *UserRepository) FindAllByRole(role model.Role) (*[]model.User, error) {
	var users []model.User
	repository.logger.Info("FindAllByRole", role)
	result := repository.database.Where("role = ?", role).Find(&users)
	if result.Error != nil {
		repository.logger.Error("Error finding users by role", result.Error)
		return nil, result.Error
	}
	return &users, nil
}
