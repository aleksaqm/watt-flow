package service

import (
	"fmt"
	"log"
	"strconv"
	"watt-flow/config"
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/repository"
	"watt-flow/util"

	"gorm.io/gorm"
)

type IElectricityConsumptionService interface {
	GetMonthlyConsumption(householdId string, year int, month int) (*dto.MonthlyConsumptionResult, error)
	Get12MonthsConsumption(householdId string, endYear int, endMonth int) (*dto.ElectricityConsumptionResponse, error)
	GetDailyConsumption(householdId string, year int, month int) (*dto.DailyConsumptionResponse, error)
	QueryConsumption(queryParams dto.FluxQueryConsumptionDto) (*dto.ConsumptionQueryResult, error)
	WithTrx(trxHandle *gorm.DB) IElectricityConsumptionService
}

type ElectricityConsumptionService struct {
	influxHelper        *util.InfluxQueryHelper
	householdRepository *repository.HouseholdRepository
}

func NewElectricityConsumptionService(env *config.Environment, householdRepository *repository.HouseholdRepository) IElectricityConsumptionService {
	influxHelper := util.NewInfluxQueryHelper(env)
	return &ElectricityConsumptionService{
		influxHelper:        influxHelper,
		householdRepository: householdRepository,
	}
}

func (s *ElectricityConsumptionService) WithTrx(trxHandle *gorm.DB) IElectricityConsumptionService {
	return &ElectricityConsumptionService{
		influxHelper:        s.influxHelper,
		householdRepository: s.householdRepository.WithTrx(trxHandle),
	}
}

func (s *ElectricityConsumptionService) GetMonthlyConsumption(householdId string, year int, month int) (*dto.MonthlyConsumptionResult, error) {
	householdIdUint, err := strconv.ParseUint(householdId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid household ID: %v", err)
	}

	household, err := s.householdRepository.FindById(householdIdUint)
	if err != nil {
		return nil, fmt.Errorf("household not found: %v", err)
	}

	if household.Status != model.OwnedHousehold {
		return nil, fmt.Errorf("household not owned")
	}

	deviceId := household.DeviceStatusID
	if deviceId == "" {
		return nil, fmt.Errorf("household has no device assigned")
	}

	consumption, err := s.influxHelper.GetTotalConsumptionForMonth(deviceId, year, month)
	if err != nil {
		log.Printf("Error getting consumption for household %s, year %d, month %d: %v", householdId, year, month, err)
		return nil, err
	}

	monthNames := []string{"", "January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December"}

	result := &dto.MonthlyConsumptionResult{
		Year:        year,
		Month:       month,
		MonthName:   monthNames[month],
		Consumption: consumption,
	}

	return result, nil
}

func (s *ElectricityConsumptionService) Get12MonthsConsumption(householdId string, endYear int, endMonth int) (*dto.ElectricityConsumptionResponse, error) {
	var results []dto.MonthlyConsumptionResult

	currentYear := endYear
	currentMonth := endMonth

	for i := 0; i < 12; i++ {
		monthResult, err := s.GetMonthlyConsumption(householdId, currentYear, currentMonth)
		if err != nil {
			log.Printf("Error getting consumption for %d/%d: %v", currentYear, currentMonth, err)
			monthNames := []string{"", "January", "February", "March", "April", "May", "June",
				"July", "August", "September", "October", "November", "December"}

			monthResult = &dto.MonthlyConsumptionResult{
				Year:        currentYear,
				Month:       currentMonth,
				MonthName:   monthNames[currentMonth],
				Consumption: 0.0,
			}
		}

		results = append(results, *monthResult)

		currentMonth--
		if currentMonth == 0 {
			currentMonth = 12
			currentYear--
		}
	}

	for i, j := 0, len(results)-1; i < j; i, j = i+1, j-1 {
		results[i], results[j] = results[j], results[i]
	}

	return &dto.ElectricityConsumptionResponse{
		Data: results,
	}, nil
}

func (s *ElectricityConsumptionService) GetDailyConsumption(householdId string, year int, month int) (*dto.DailyConsumptionResponse, error) {
	householdIdUint, err := strconv.ParseUint(householdId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid household ID: %v", err)
	}

	household, err := s.householdRepository.FindById(householdIdUint)
	if err != nil {
		return nil, fmt.Errorf("household not found: %v", err)
	}

	if household.Status != model.OwnedHousehold {
		return nil, fmt.Errorf("household not owned")
	}

	deviceId := household.DeviceStatusID
	if deviceId == "" {
		return nil, fmt.Errorf("household has no device assigned")
	}

	daysInMonth := getDaysInMonth(year, month)

	var results []dto.DailyConsumptionData

	for day := 1; day <= daysInMonth; day++ {
		consumption, err := s.influxHelper.GetTotalConsumptionForDay(deviceId, year, month, day)
		if err != nil {
			log.Printf("Warning: Failed to get consumption for day %d/%d/%d: %v", day, month, year, err)
			consumption = 0.0
		}

		results = append(results, dto.DailyConsumptionData{
			Year:        year,
			Month:       month,
			Day:         day,
			Consumption: consumption,
		})
	}

	return &dto.DailyConsumptionResponse{
		Data: results,
	}, nil
}

func (s *ElectricityConsumptionService) QueryConsumption(queryParams dto.FluxQueryConsumptionDto) (*dto.ConsumptionQueryResult, error) {
	result, err := s.influxHelper.SendConsumptionQuery(queryParams)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getDaysInMonth(year, month int) int {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if (year%4 == 0 && year%100 != 0) || (year%400 == 0) {
			return 29
		}
		return 28
	default:
		return 30
	}
}
