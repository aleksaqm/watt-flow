package service

import (
	"time"
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/repository"
)

type IPropertyService interface {
	FindById(id uint64) (*model.Property, error)
	Create(property *dto.CreatePropertyDto) (*model.Property, error)
	Update(property *model.Property) (*model.Property, error)
	Delete(id uint64) error
	FindByStatus(status model.PropertyStatus) ([]model.Property, error)
}

type PropertyService struct {
	repository       *repository.PropertyRepository
	householdService IHouseholdService
}

func NewPropertyService(repository *repository.PropertyRepository, householdService IHouseholdService) *PropertyService {
	return &PropertyService{
		repository:       repository,
		householdService: householdService,
	}
}

func (service *PropertyService) FindById(id uint64) (*model.Property, error) {
	property, err := service.repository.FindById(id)
	if err != nil {
		return nil, err
	}
	return property, nil
}

func (service *PropertyService) FindByStatus(status model.PropertyStatus) ([]model.Property, error) {
	properties, err := service.repository.FindByStatus(status)
	if err != nil {
		return nil, err
	}
	return properties, nil
}

func (service *PropertyService) Create(propertyDto *dto.CreatePropertyDto) (*model.Property, error) {
	property := model.Property{}
	property.CreatedAt = time.Now()
	property.Images = propertyDto.Images
	property.Documents = propertyDto.Documents
	property.AddressID = propertyDto.AddressID
	property.Floors = propertyDto.Floors

	// TESTING
	property.OwnerID = 1

	createdProperty, err := service.repository.Create(&property)
	if err != nil {
		return nil, err
	}
	return &createdProperty, nil
}

func (s *PropertyService) SearchHouseholds(searchDto dto.SearchHouseholdDto) ([]model.Household, error) {
	var households []model.Household
	if searchDto.Id != "" {
		household, err := s.householdService.FindByCadastralNumber(searchDto.Id)
		if err != nil {
			return nil, err
		}
		households = make([]model.Household, 0)
		households = append(households, *household)
		return households, nil
	}
	households, err := s.householdService.Search(searchDto)
	if err != nil {
		return nil, err
	}
	return households, nil
}

func (service *PropertyService) Update(property *model.Property) (*model.Property, error) {
	updatedProperty, err := service.repository.Update(property)
	if err != nil {
		return nil, err
	}
	return &updatedProperty, nil
}

func (service *PropertyService) Delete(id uint64) error {
	err := service.repository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
