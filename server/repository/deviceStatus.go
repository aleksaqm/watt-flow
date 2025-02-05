package repository

import (
	"watt-flow/db"
	"watt-flow/model"
	"watt-flow/util"
)

type DeviceStatusRepository struct {
	database db.Database
	logger   util.Logger
}

func NewDeviceStatusRepository(db db.Database, logger util.Logger) DeviceStatusRepository {
	err := db.AutoMigrate(&model.DeviceStatus{})
	if err != nil {
		logger.Error("Error migrating device status", err)
	}
	return DeviceStatusRepository{
		database: db,
		logger:   logger,
	}
}

func (repository *DeviceStatusRepository) Create(deviceStatus *model.DeviceStatus) (model.DeviceStatus, error) {
	result := repository.database.Create(deviceStatus)
	if result.Error != nil {
		repository.logger.Error("Error creating device status", result.Error)
		return *deviceStatus, result.Error
	}
	return *deviceStatus, nil
}

func (repository *DeviceStatusRepository) FindById(address string) (*model.DeviceStatus, error) {
	var deviceStatus model.DeviceStatus
	if err := repository.database.First(&deviceStatus, "address = ?", address).Error; err != nil {
		repository.logger.Error("Error finding device status by address", err)
		return nil, err
	}
	return &deviceStatus, nil
}

func (repository *DeviceStatusRepository) FindByHouseholdID(householdID uint64) (*model.DeviceStatus, error) {
	var deviceStatus model.DeviceStatus
	if err := repository.database.Where("household_id = ?", householdID).First(&deviceStatus).Error; err != nil {
		repository.logger.Error("Error finding device status by household ID", err)
		return nil, err
	}
	return &deviceStatus, nil
}

func (repository *DeviceStatusRepository) Update(deviceStatus *model.DeviceStatus) (model.DeviceStatus, error) {
	result := repository.database.Save(deviceStatus)
	if result.Error != nil {
		repository.logger.Error("Error updating device status", result.Error)
		return *deviceStatus, result.Error
	}
	return *deviceStatus, nil
}

func (repository *DeviceStatusRepository) Delete(address string) error {
	result := repository.database.Delete(&model.DeviceStatus{}, "address = ?", address)
	if result.Error != nil {
		repository.logger.Error("Error deleting device status", result.Error)
		return result.Error
	}
	return nil
}
