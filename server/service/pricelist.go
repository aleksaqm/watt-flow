package service

import (
	"watt-flow/model"
	"watt-flow/repository"
)

type IPricelistService interface {
	CreatePricelist(pricelist *model.Pricelist) (*model.Pricelist, error)
	FindById(id uint64) (*model.Pricelist, error)
}

type PricelistService struct {
	pricelistRepository *repository.PricelistRepository
}

func NewPricelistService(pricelistRepository *repository.PricelistRepository) *PricelistService {
	return &PricelistService{
		pricelistRepository: pricelistRepository,
	}
}

func (t *PricelistService) FindById(id uint64) (*model.Pricelist, error) {
	pricelist, err := t.pricelistRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	return pricelist, nil
}

func (t *PricelistService) CreatePricelist(pricelist *model.Pricelist) (*model.Pricelist, error) {
	_, err := t.pricelistRepository.Create(pricelist)
	if err != nil {
		return nil, err
	}
	return pricelist, nil
}
