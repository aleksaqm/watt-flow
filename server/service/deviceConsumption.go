package service

import (
	"watt-flow/dto"
	"watt-flow/util"
)

type IDeviceConsumptionService interface {
	QueryConsumption(queryParams dto.FluxQueryConsumptionDto) (*dto.ConsumptionQueryResult, error)
}

type DeviceConsumptionService struct {
	influxQueryHelper *util.InfluxQueryHelper
}

func NewDeviceConsumptionService(influx *util.InfluxQueryHelper) *DeviceConsumptionService {
	return &DeviceConsumptionService{
		influxQueryHelper: influx,
	}
}

func (service *DeviceConsumptionService) QueryConsumption(queryParams dto.FluxQueryConsumptionDto) (*dto.ConsumptionQueryResult, error) {
	result, err := service.influxQueryHelper.SendConsumptionQuery(queryParams)
	if err != nil {
		return nil, err
	}
	return result, nil
}
