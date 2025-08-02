package service

import (
	"watt-flow/repository"

	"gorm.io/gorm"
)

type ICityService interface {
	GetAllCities() ([]string, error)
	WithTrx(trxHandle *gorm.DB) ICityService
}

type CityService struct {
	repository repository.CityRepository
}

func NewCityService(repository repository.CityRepository) *CityService {
	return &CityService{
		repository: repository,
	}
}

func (s CityService) WithTrx(trxHandle *gorm.DB) ICityService {
	s.repository = s.repository.WithTrx(trxHandle)
	return &s
}

func (service *CityService) GetAllCities() ([]string, error) {
	return service.repository.GetAllCities()
}
