package service

import (
	"watt-flow/model"
	"watt-flow/repository"
)

type IPropertyService interface {
	FindById(id uint64) (*model.Property, error)
	Create(property *model.Property) (*model.Property, error)
	Update(property *model.Property) (*model.Property, error)
	Delete(id uint64) error
	FindByStatus(status model.PropertyStatus) ([]model.Property, error)
}

type PropertyService struct {
	repository *repository.PropertyRepository
}

func NewPropertyService(repository *repository.PropertyRepository) *PropertyService {
	return &PropertyService{
		repository: repository,
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

func (service *PropertyService) Create(property *model.Property) (*model.Property, error) {
	createdProperty, err := service.repository.Create(property)
	if err != nil {
		return nil, err
	}
	return &createdProperty, nil
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
