package service

import (
	"fmt"
	"time"
	"watt-flow/model"
	"watt-flow/repository"
)

type IBillService interface {
	FindById(id uint64) (*model.Bill, error)
	// Query(params *dto.BillQueryParams) ([]model.Bill, int64, error)
	GetUnsentMonthlyBills() ([]string, error)
}

type BillService struct {
	billRepository        *repository.BillRepository
	monthlyBillRepository *repository.MonthlyBillRepository
}

func NewBillService(billRepository *repository.BillRepository, monthlyBillRepository *repository.MonthlyBillRepository) *BillService {
	return &BillService{
		billRepository:        billRepository,
		monthlyBillRepository: monthlyBillRepository,
	}
}

func (t *BillService) FindById(id uint64) (*model.Bill, error) {
	bill, err := t.billRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	return bill, nil
}

func formatMonthKey(month, year uint) string {
	return fmt.Sprintf("%d-%02d", year, month)
}

func isValidBillingMonth(date time.Time) bool {
	now := time.Now()
	if date.Year() > now.Year() {
		return false
	}
	if date.Year() == now.Year() && date.Month() >= now.Month() {
		return false
	}

	endOfMonth := time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	// buffer period of one day
	bufferPeriod := endOfMonth.AddDate(0, 0, 1)

	return now.After(bufferPeriod)
}

func (s *BillService) GetUnsentMonthlyBills() ([]string, error) {
	var missingMonths []string

	endDate := time.Now()
	startDate := endDate.AddDate(-1, 0, 0)

	existingBills, err := s.monthlyBillRepository.FindBillsBetweenDates(startDate, endDate)
	if err != nil {
		return nil, err
	}

	existingBillsMap := make(map[string]bool)
	for _, bill := range existingBills {
		existingBillsMap[bill.BillingDate] = true
	}

	currentDate := startDate
	for currentDate.Before(endDate) {
		key := fmt.Sprintf("%d-%02d", currentDate.Year(), int(currentDate.Month()))
		if !existingBillsMap[key] && isValidBillingMonth(currentDate) {
			missingMonths = append(missingMonths, key)
		}
		currentDate = currentDate.AddDate(0, 1, 0)
	}

	return missingMonths, nil
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
