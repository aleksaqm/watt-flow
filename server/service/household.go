package service

import (
	"fmt"
	"log"
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/repository"

	"gorm.io/gorm"
)

type IHouseholdService interface {
	FindById(id uint64) (*model.Household, error)
	Create(household *dto.CreateHouseholdDto) (*model.Household, error)
	Update(household *model.Household) (*model.Household, error)
	Delete(id uint64) error
	FindByStatus(status model.HouseholdStatus) ([]model.Household, error)
	FindByCadastralNumber(id string) (*model.Household, error)
	Query(queryParams *dto.HouseholdQueryParams) ([]dto.HouseholdResultDto, int64, error)
	AcceptHouseholds(tx *gorm.DB, propertyID uint64) error
}

type HouseholdService struct {
	repository *repository.HouseholdRepository
}

func NewHouseholdService(repository *repository.HouseholdRepository) *HouseholdService {
	return &HouseholdService{
		repository: repository,
	}
}

func (service *HouseholdService) FindById(id uint64) (*model.Household, error) {
	household, err := service.repository.FindById(id)
	if err != nil {
		return nil, err
	}
	return household, nil
}

func (service *HouseholdService) Query(queryParams *dto.HouseholdQueryParams) ([]dto.HouseholdResultDto, int64, error) {
	var households []dto.HouseholdResultDto
	if queryParams.Search.Id != "" {
		household, err := service.FindByCadastralNumber(queryParams.Search.Id)
		fmt.Printf("household: %v", household)
		if err != nil {
			return nil, 0, err
		}
		households = make([]dto.HouseholdResultDto, 0)
		if household != nil {
			mapped_household, _ := MapToResultDto(household)

			households = append(households, mapped_household)

			fmt.Printf("household: %v", households)
			return households, 1, nil
		}

		return households, 0, nil
	}

	data, count, err := service.repository.Query(queryParams)
	if err != nil {
		log.Printf("Error on query: %v", err)
		return nil, 0, err
	}
	households = make([]dto.HouseholdResultDto, 0)
	for _, household := range data {
		mapped_household, _ := MapToResultDto(&household)
		households = append(households, mapped_household)
	}

	return households, count, nil
}

func MapToResultDto(household *model.Household) (dto.HouseholdResultDto, error) {
	ownerUsername := ""
	var ownerId uint64 = 0
	if household.OwnerID != nil {
		ownerUsername = household.Owner.Username
		ownerId = *household.OwnerID
	}
	response := dto.HouseholdResultDto{
		Id:              household.Id,
		Floor:           household.Floor,
		Suite:           household.Suite,
		Status:          household.Status.HouseholdStatusToString(),
		SqFootage:       household.SqFootage,
		OwnerID:         ownerId,
		OwnerName:       ownerUsername,
		MeterAddress:    household.DeviceStatus.DeviceId,
		PropertyID:      household.PropertyID,
		CadastralNumber: household.CadastralNumber,
		City:            household.Property.Address.City,
		Street:          household.Property.Address.Street,
		Number:          household.Property.Address.Number,
	}
	return response, nil
}

func (service *HouseholdService) FindByCadastralNumber(id string) (*model.Household, error) {
	household, err := service.repository.FindByCadastralNumber(id)
	if err != nil {
		return nil, err
	}
	return household, nil
}

func (service *HouseholdService) FindByStatus(status model.HouseholdStatus) ([]model.Household, error) {
	households, err := service.repository.FindByStatus(status)
	if err != nil {
		return nil, err
	}
	return households, nil
}

func (service *HouseholdService) Create(householdDto *dto.CreateHouseholdDto) (*model.Household, error) {
	household := model.Household{
		// Initialize fields from the DTO
		Floor:     householdDto.Floor,
		Suite:     householdDto.Suite,
		SqFootage: householdDto.SqFootage,
		Status:    model.InactiveHousehold,
		// OwnerID:         1,
		PropertyID:      householdDto.PropertyId,
		CadastralNumber: householdDto.CadastralNumber,
	}

	createdHousehold, err := service.repository.Create(&household)
	if err != nil {
		return nil, err
	}
	return &createdHousehold, nil
}

func (service *HouseholdService) Update(household *model.Household) (*model.Household, error) {
	updatedHousehold, err := service.repository.Update(household)
	if err != nil {
		return nil, err
	}
	return &updatedHousehold, nil
}

func (service *HouseholdService) Delete(id uint64) error {
	err := service.repository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (service *HouseholdService) AcceptHouseholds(tx *gorm.DB, propertyID uint64) error {
	err := service.repository.AcceptHouseholds(tx, propertyID)
	if err != nil {
		return err
	}
	return nil
}
