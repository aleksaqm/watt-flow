package service

import (
	"watt-flow/repository"
)

type ICityService interface {
	GetAllCities() ([]string, error)
}

type CityService struct {
	repository repository.CityRepository
}

func NewCityService(repository repository.CityRepository) *CityService {
	return &CityService{
		repository: repository,
	}
}

func (service *CityService) GetAllCities() ([]string, error) {
	return service.repository.GetAllCities()
}
