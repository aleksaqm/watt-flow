package repository

import (
	"fmt"

	"watt-flow/db"
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/util"

	"gorm.io/gorm"
)

type UserRepository struct {
	Database db.Database
	Logger   util.Logger
}

func NewUserRepository(db db.Database, logger util.Logger) UserRepository {
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		logger.Error("Error migrating user", err)
	}
	return UserRepository{
		Database: db,
		Logger:   logger,
	}
}

func (r UserRepository) WithTrx(trxHandle *gorm.DB) UserRepository {
	if trxHandle == nil {
		r.Logger.Error("Transaction Database not found in gin context. ")
		return r
	}
	r.Database.DB = trxHandle
	return r
}

func (repository *UserRepository) Create(user *model.User) (model.User, error) {
	result := repository.Database.Create(user)
	if result.Error != nil {
		repository.Logger.Error("Error creating user", result.Error)
		return *user, result.Error
	}
	return *user, nil
}

func (repository *UserRepository) FindActiveByEmail(email string) (*model.User, error) {
	var user model.User
	result := repository.Database.Where("email = ? AND status = ?", email, 0).First(&user)
	if result.Error != nil {
		repository.Logger.Error("Error finding user by email", result.Error)
		return nil, result.Error
	}
	return &user, nil
}

func (repository *UserRepository) FindById(id uint64) (*model.User, error) {
	var user model.User
	if err := repository.Database.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repository *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	result := repository.Database.Where("username = ?", username).First(&user)
	if result.Error != nil {
		repository.Logger.Error("Error finding user by username", result.Error)
		return nil, result.Error
	}
	return &user, nil
}

func (repository *UserRepository) FindByUsernameOrActiveEmail(username string, email string) (*model.User, error) {
	var user model.User
	result := repository.Database.Where("username = ? OR (email = ? AND status = ?)", username, email, 0).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repository *UserRepository) FindByEmailAndUsername(email string, username string) (*model.User, error) {
	var user model.User
	result := repository.Database.Where("email = ? AND username = ?", email, username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repository *UserRepository) Query(params *dto.UserQueryParams) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	baseQuery := repository.Database.Model(&model.User{})
	role, _ := model.ParseRole(params.Search.Role)

	if params.Search.Role != "" {
		baseQuery = baseQuery.Where("role = ?", role)
	}
	if params.Search.Status != "" {
		status, err := model.ParseAccountStatus(params.Search.Status)
		if err != nil {
			repository.Logger.Error("Error parsing account status", err)
			return nil, 0, err
		}
		baseQuery = baseQuery.Where("status = ?", status)
	}
	if params.Search.Username != "" {
		baseQuery = baseQuery.Where("username ILIKE ?", "%"+params.Search.Username+"%")
	}
	// if params.Search.Number != "" {
	// 	baseQuery = baseQuery.Where("number ILIKE ?", "%"+params.Search.Number+"%")
	// }
	if err := baseQuery.Count(&total).Error; err != nil {
		repository.Logger.Error("Error querying user count", err)
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
		repository.Logger.Error("Error querying users", err)
		return nil, 0, err
	}

	return users, total, nil
}

func (repository *UserRepository) Update(user *model.User) (model.User, error) {
	result := repository.Database.Save(user)
	if result.Error != nil {
		repository.Logger.Error("Error updating user", result.Error)
		return *user, result.Error
	}
	return *user, nil
}

func (repository *UserRepository) FindAllByRole(role model.Role) (*[]model.User, error) {
	var users []model.User
	repository.Logger.Info("FindAllByRole", role)
	result := repository.Database.Where("role = ?", role).Find(&users)
	if result.Error != nil {
		repository.Logger.Error("Error finding users by role", result.Error)
		return nil, result.Error
	}
	return &users, nil
}

func (repository *UserRepository) ChangeStatus(userId uint64, status int) error {
	err := repository.Database.Model(&model.User{}).
		Where("id = ?", userId).
		Update("status", status).
		Error
	if err != nil {
		repository.Logger.Error("Error updating account status", err)
		return err
	}
	return nil
}
