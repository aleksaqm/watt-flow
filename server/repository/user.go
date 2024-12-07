package repository

import (
	"fmt"
	"watt-flow/db"
	"watt-flow/dto"
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

func (repository *UserRepository) FindActiveByEmail(email string) (*model.User, error) {
	var user model.User
	result := repository.database.Where("email = ? AND status = ?", email, 0).First(&user)
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

func (repository *UserRepository) FindByUsernameOrActiveEmail(username string, email string) (*model.User, error) {
	var user model.User
	result := repository.database.Where("username = ? OR (email = ? AND status = ?)", username, email, 0).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repository *UserRepository) FindByEmailAndUsername(email string, username string) (*model.User, error) {
	var user model.User
	result := repository.database.Where("email = ? AND username = ?", email, username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repository *UserRepository) Query(params *dto.UserQueryParams) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	baseQuery := repository.database.Model(&model.User{})
	role, _ := model.ParseRole(params.Search.Role)

	if params.Search.Role != "" {
		baseQuery = baseQuery.Where("role = ?", role)
	}
	if params.Search.Username != "" {
		baseQuery = baseQuery.Where("username ILIKE ?", "%"+params.Search.Username+"%")
	}
	// if params.Search.Number != "" {
	// 	baseQuery = baseQuery.Where("number ILIKE ?", "%"+params.Search.Number+"%")
	// }
	if err := baseQuery.Count(&total).Error; err != nil {
		repository.logger.Error("Error querying user count", err)
		return nil, 0, err
	}

	sortOrder := params.SortOrder
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}
	query := baseQuery.Order(fmt.Sprintf("%s %s", params.SortBy, sortOrder))
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	if err := query.
		Find(&users).Error; err != nil {
		repository.logger.Error("Error querying users", err)
		return nil, 0, err
	}

	return users, total, nil
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

func (repository *UserRepository) ChangeStatus(userId uint64, status int) error {
	err := repository.database.Model(&model.User{}).
		Where("id = ?", userId).
		Update("status", status).
		Error
	if err != nil {
		repository.logger.Error("Error updating account status", err)
		return err
	}
	return nil
}
