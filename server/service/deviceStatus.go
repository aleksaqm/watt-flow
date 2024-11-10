package service

import (
	"watt-flow/model"
	"watt-flow/repository"
)

type IDeviceStatusService interface {
	FindByAddress(address string) (*model.DeviceStatus, error)
	FindByHouseholdID(householdID uint64) (*model.DeviceStatus, error)
	Create(deviceStatus *model.DeviceStatus) (*model.DeviceStatus, error)
	Update(deviceStatus *model.DeviceStatus) (*model.DeviceStatus, error)
	Delete(address string) error
}

type DeviceStatusService struct {
	repository *repository.DeviceStatusRepository
}

func NewDeviceStatusService(repository *repository.DeviceStatusRepository) *DeviceStatusService {
	return &DeviceStatusService{
		repository: repository,
	}
}

func (service *DeviceStatusService) FindByAddress(address string) (*model.DeviceStatus, error) {
	deviceStatus, err := service.repository.FindById(address)
	if err != nil {
		return nil, err
	}
	return deviceStatus, nil
}

func (service *DeviceStatusService) FindByHouseholdID(householdID uint64) (*model.DeviceStatus, error) {
	deviceStatus, err := service.repository.FindByHouseholdID(householdID)
	if err != nil {
		return nil, err
	}
	return deviceStatus, nil
}

func (service *DeviceStatusService) Create(deviceStatus *model.DeviceStatus) (*model.DeviceStatus, error) {
	createdDeviceStatus, err := service.repository.Create(deviceStatus)
	if err != nil {
		return nil, err
	}
	return &createdDeviceStatus, nil
}

func (service *DeviceStatusService) Update(deviceStatus *model.DeviceStatus) (*model.DeviceStatus, error) {
	updatedDeviceStatus, err := service.repository.Update(deviceStatus)
	if err != nil {
		return nil, err
	}
	return &updatedDeviceStatus, nil
}

func (service *DeviceStatusService) Delete(address string) error {
	err := service.repository.Delete(address)
	if err != nil {
		return err
	}
	return nil
}
