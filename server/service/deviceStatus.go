package service

import (
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/repository"
	"watt-flow/util"
)

type IDeviceStatusService interface {
	FindByAddress(address string) (*model.DeviceStatus, error)
	FindByHouseholdID(householdID uint64) (*model.DeviceStatus, error)
	Create(deviceStatus *model.DeviceStatus) (*model.DeviceStatus, error)
	Update(deviceStatus *model.DeviceStatus) (*model.DeviceStatus, error)
	Delete(address string) error
	QueryStatus(queryParams dto.FluxQueryStatusDto) (*dto.StatusQueryResult, error)
	QueryConsumption(queryParams dto.FluxQueryConsumptionDto) (*dto.StatusQueryResult, error)
}

type DeviceStatusService struct {
	repository        repository.DeviceStatusRepository
	influxQueryHelper *util.InfluxQueryHelper
}

func NewDeviceStatusService(repository repository.DeviceStatusRepository, influx *util.InfluxQueryHelper) *DeviceStatusService {
	return &DeviceStatusService{
		repository:        repository,
		influxQueryHelper: influx,
	}
}

func (service *DeviceStatusService) QueryStatus(queryParams dto.FluxQueryStatusDto) (*dto.StatusQueryResult, error) {
	result, err := service.influxQueryHelper.SendStatusQuery(queryParams)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (service *DeviceStatusService) QueryConsumption(queryParams dto.FluxQueryConsumptionDto) (*dto.StatusQueryResult, error) {
	result, err := service.influxQueryHelper.SendConsumptionQuery(queryParams)
	if err != nil {
		return nil, err
	}
	return result, nil
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
