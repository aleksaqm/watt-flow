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
	repository *repository.CityRepository
}

func NewCityService(repository *repository.CityRepository) *CityService {
	return &CityService{
		repository: repository,
	}
}

func (s *CityService) WithTrx(trxHandle *gorm.DB) ICityService {
	return &CityService{
		repository: s.repository.WithTrx(trxHandle),
	}
}

func (service *CityService) GetAllCities() ([]string, error) {
	return service.repository.GetAllCities()
}
