package repository

import (
	"watt-flow/db"
	"watt-flow/model"
	"watt-flow/util"

	"gorm.io/gorm"
)

type AddressRepository struct {
	database db.Database
	logger   util.Logger
}

func NewAddressRepository(db db.Database, logger util.Logger) AddressRepository {
	err := db.AutoMigrate(&model.Address{})
	if err != nil {
		logger.Error("Error migrating address", err)
	}
	return AddressRepository{
		database: db,
		logger:   logger,
	}
}

func (r AddressRepository) WithTrx(trxHandle *gorm.DB) AddressRepository {
	if trxHandle == nil {
		r.logger.Error("Transaction Database not found in gin context. ")
		return r
	}
	r.database.DB = trxHandle
	return r
}

func (repository *AddressRepository) Create(address *model.Address) (*model.Address, error) {
	result := repository.database.Create(address)
	if result.Error != nil {
		repository.logger.Error("Error creating address", result.Error)
		return nil, result.Error
	}
	return address, nil
}

func (repository *AddressRepository) FindById(id uint64) (*model.Address, error) {
	var address model.Address
	if err := repository.database.First(&address, id).Error; err != nil {
		repository.logger.Error("Error finding address by ID", err)
		return nil, err
	}
	return &address, nil
}

func (repository *AddressRepository) FindAll() ([]model.Address, error) {
	var addresses []model.Address
	result := repository.database.Find(&addresses)
	if result.Error != nil {
		repository.logger.Error("Error finding all addresses", result.Error)
		return nil, result.Error
	}
	return addresses, nil
}

func (repository *AddressRepository) Update(address *model.Address) (*model.Address, error) {
	result := repository.database.Save(address)
	if result.Error != nil {
		repository.logger.Error("Error updating address", result.Error)
		return nil, result.Error
	}
	return address, nil
}

func (repository *AddressRepository) Delete(id uint64) error {
	result := repository.database.Delete(&model.Address{}, id)
	if result.Error != nil {
		repository.logger.Error("Error deleting address", result.Error)
		return result.Error
	}
	return nil
}
