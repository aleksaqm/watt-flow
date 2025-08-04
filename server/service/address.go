package service

import (
	"watt-flow/model"
	"watt-flow/repository"

	"gorm.io/gorm"
)

type IAddressService interface {
	Create(address *model.Address) (*model.Address, error)
	FindById(id uint64) (*model.Address, error)
	FindAll() ([]model.Address, error)
	Update(address *model.Address) (*model.Address, error)
	Delete(id uint64) error
	WithTrx(trxHandle *gorm.DB) IAddressService
}

type AddressService struct {
	repository *repository.AddressRepository
}

func NewAddressService(repository *repository.AddressRepository) *AddressService {
	return &AddressService{
		repository: repository,
	}
}

func (s *AddressService) WithTrx(trxHandle *gorm.DB) IAddressService {
	return &AddressService{
		repository: s.repository.WithTrx(trxHandle),
	}
}

func (service *AddressService) Create(address *model.Address) (*model.Address, error) {
	createdAddress, err := service.repository.Create(address)
	if err != nil {
		return nil, err
	}
	return createdAddress, nil
}

func (service *AddressService) FindById(id uint64) (*model.Address, error) {
	address, err := service.repository.FindById(id)
	if err != nil {
		return nil, err
	}
	return address, nil
}

func (service *AddressService) FindAll() ([]model.Address, error) {
	addresses, err := service.repository.FindAll()
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (service *AddressService) Update(address *model.Address) (*model.Address, error) {
	updatedAddress, err := service.repository.Update(address)
	if err != nil {
		return nil, err
	}
	return updatedAddress, nil
}

func (service *AddressService) Delete(id uint64) error {
	err := service.repository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
