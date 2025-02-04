package service

import (
	"fmt"
	"time"
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/repository"
	"watt-flow/util"
)

type IBillService interface {
	FindById(id uint64) (*model.Bill, error)
	GenerateMonthlyBill(year int, month int) (*model.MonthlyBill, error)
	QueryMonthly(params *dto.MonthlyBillQueryParams) ([]model.MonthlyBill, int64, error)
	GetUnsentMonthlyBills() ([]string, error)
	InitiateBilling(year int, month int) (*model.MonthlyBill, error)
}

type BillService struct {
	billRepository        *repository.BillRepository
	monthlyBillRepository *repository.MonthlyBillRepository
	householdService      IHouseholdService
	pricelistService      IPricelistService
	influxQueryHelper     *util.InfluxQueryHelper
	emailSender           *util.EmailSender
}

func NewBillService(billRepository *repository.BillRepository, monthlyBillRepository *repository.MonthlyBillRepository, householdService IHouseholdService, pricelistService IPricelistService, influxQueryHelper *util.InfluxQueryHelper, emailSender *util.EmailSender) *BillService {
	return &BillService{
		billRepository:        billRepository,
		monthlyBillRepository: monthlyBillRepository,
		householdService:      householdService,
		pricelistService:      pricelistService,
		influxQueryHelper:     influxQueryHelper,
		emailSender:           emailSender,
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

func (service *BillService) QueryMonthly(queryParams *dto.MonthlyBillQueryParams) ([]model.MonthlyBill, int64, error) {
	var bills []model.MonthlyBill
	bills, total, err := service.monthlyBillRepository.Query(queryParams)
	if err != nil {
		return nil, 0, err
	}
	return bills, total, nil
}

func (service *BillService) InitiateBilling(year int, month int) (*model.MonthlyBill, error) {
	monthlyBill, err := service.GenerateMonthlyBill(year, month)
	if err != nil {
		return nil, err
	}
	households, err := service.householdService.GetOwnedHouseholds()
	if err != nil {
		return nil, err
	}
	activePricelist, err := service.pricelistService.GetActivePricelist()
	if err != nil {
		return nil, err
	}
	for _, household := range households {
		spentPower, err := service.influxQueryHelper.GetTotalConsumptionForMonth(household.DeviceStatusID, year, month)
		if err != nil {
			return nil, err
		}
		calculatedPrice := calculatePrice(spentPower, *activePricelist)
		if household.Owner == nil {
			fmt.Println("Error owner")
		}

		if activePricelist == nil {
			fmt.Println("Error pricelist")
		}
		billingDate := fmt.Sprintf("%d-%02d", year, month)

		bill := &model.Bill{
			BillingDate: billingDate,
			IssueDate:   time.Now(),
			Pricelist:   *activePricelist,
			Owner:       *household.Owner,
			SpentPower:  spentPower,
			Price:       calculatedPrice,
			PricelistID: activePricelist.ID,
			OwnerID:     household.Owner.Id,
		}
		fmt.Println(bill)
		emailBody := util.GenerateMonthlyBillEmail(bill)
		err = service.emailSender.SendEmail(household.Owner.Email, "Electricity bill for "+billingDate, emailBody)
		if err != nil {
			return nil, err
		}
		_, err = service.billRepository.Create(bill)
		fmt.Println("saved bill")
		if err != nil {
			return nil, err
		}
	}
	return monthlyBill, nil
}

func calculatePrice(spentPower float64, pricelist model.Pricelist) float64 {
	const billingPowerConstant = 7.0

	greenConsumption := min(spentPower, 350)
	blueConsumption := min(max(spentPower-350, 0), 1250)
	redConsumption := max(spentPower-1600, 0)

	basePrice := (greenConsumption * pricelist.GreenZone) +
		(blueConsumption * pricelist.BlueZone) +
		(redConsumption * pricelist.RedZone) +
		(billingPowerConstant * pricelist.BillingPower)

	// Apply tax
	finalPrice := basePrice + (basePrice * pricelist.Tax / 100)

	return finalPrice
}

func (s *BillService) GenerateMonthlyBill(year int, month int) (*model.MonthlyBill, error) {
	now := time.Now()

	billingDate, err := time.Parse("2006-01", fmt.Sprintf("%d-%02d", year, month))
	if err != nil {
		return nil, fmt.Errorf("invalid year or month format")
	}

	// Check if the billing date is in the future
	if billingDate.After(time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)) {
		return nil, fmt.Errorf("cannot issue a bill for a future month")
	}

	monthlyBill := model.MonthlyBill{
		IssueDate:   time.Now(),
		BillingDate: fmt.Sprintf("%d-%02d", year, month),
		Status:      "In Progress",
	}
	bill, err := s.monthlyBillRepository.Create(&monthlyBill)
	if err != nil {
		return nil, err
	}
	return &bill, nil
}
