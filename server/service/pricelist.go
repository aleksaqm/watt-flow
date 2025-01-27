package service

import (
	"fmt"
	"time"
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/repository"

	"gorm.io/datatypes"
)

type IPricelistService interface {
	CreatePricelist(pricelist *dto.NewPricelist) (*model.Pricelist, error)
	FindById(id uint64) (*model.Pricelist, error)
	FindByDate(date datatypes.Date) (*model.Pricelist, error)
	Query(params *dto.PricelistQueryParams) ([]model.Pricelist, int64, error)
	Delete(id uint64) error
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

func (t *PricelistService) FindByDate(date datatypes.Date) (*model.Pricelist, error) {
	pricelist, err := t.pricelistRepository.FindByDate(date)
	if err != nil {
		return nil, err
	}
	return pricelist, nil
}

func (t *PricelistService) Delete(id uint64) error {
	pricelist, err := t.pricelistRepository.FindById(id)
	if err != nil {
		return err
	}
	if time.Time(pricelist.ValidFrom).Before(time.Now().UTC()) {
		return fmt.Errorf("failed to delete pricelist in the past")
	}

	return t.pricelistRepository.Delete(id)
}

func (t *PricelistService) CreatePricelist(newPricelist *dto.NewPricelist) (*model.Pricelist, error) {
	date := datatypes.Date(time.Date(newPricelist.Year, time.Month(newPricelist.Month), 1, 0, 0, 0, 0, time.UTC))
	found, _ := t.FindByDate(date)
	if found != nil {
		return nil, fmt.Errorf("pricelist already exists for given time")
	}
	pricelist := model.Pricelist{
		ValidFrom:    date,
		RedZone:      newPricelist.RedZone,
		BlueZone:     newPricelist.BlueZone,
		GreenZone:    newPricelist.GreenZone,
		BillingPower: newPricelist.BillingPower,
		Tax:          newPricelist.Tax,
	}
	createdPricelist, err := t.pricelistRepository.Create(&pricelist)
	if err != nil {
		return nil, err
	}
	return &createdPricelist, nil
}

func (t *PricelistService) Query(params *dto.PricelistQueryParams) ([]model.Pricelist, int64, error) {
	var pricelists []model.Pricelist
	pricelists, total, err := t.pricelistRepository.Query(params)
	return pricelists, total, err
}
