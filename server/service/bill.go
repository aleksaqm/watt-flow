package service

import (
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/repository"
)

type IBillService interface {
	FindById(id uint64) (*model.Bill, error)
	Query(params *dto.BillQueryParams) ([]model.Bill, int64, error)
}

type BillService struct {
	billRepository *repository.BillRepository
}

func NewBillService(billRepository *repository.BillRepository) *BillService {
	return &BillService{
		billRepository: billRepository,
	}
}

func (t *BillService) FindById(id uint64) (*model.Bill, error) {
	bill, err := t.billRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	return bill, nil
}

// func (t *BillService) CreatePricelist(newPricelist *dto.NewPricelist) (*model.Pricelist, error) {
// 	date := datatypes.Date(time.Date(newPricelist.Year, time.Month(newPricelist.Month), 1, 0, 0, 0, 0, time.UTC))
// 	found, _ := t.FindByDate(date)
// 	if found != nil {
// 		return nil, fmt.Errorf("pricelist already exists for given time")
// 	}
// 	pricelist := model.Pricelist{
// 		ValidFrom:    date,
// 		RedZone:      newPricelist.RedZone,
// 		BlueZone:     newPricelist.BlueZone,
// 		GreenZone:    newPricelist.GreenZone,
// 		BillingPower: newPricelist.BillingPower,
// 		Tax:          newPricelist.Tax,
// 	}
// 	createdPricelist, err := t.pricelistRepository.Create(&pricelist)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &createdPricelist, nil
// }

// func (t *BillService) Query(params *dto.BillQueryParams) ([]model.Bill, int64, error) {
// 	var bills []model.Bill
// 	bills, total, err := t.billRepository.Query(params)
// 	return bills, total, err
// }
