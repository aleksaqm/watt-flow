package service

import (
	"gorm.io/gorm"
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/repository"
)

type IHouseholdService interface {
	FindById(id uint64) (*model.Household, error)
	Create(household *dto.CreateHouseholdDto) (*model.Household, error)
	Update(household *model.Household) (*model.Household, error)
	Delete(id uint64) error
	FindByStatus(status model.HouseholdStatus) ([]model.Household, error)
	FindByCadastralNumber(id string) (*model.Household, error)
	Search(searchDto dto.SearchHouseholdDto) ([]model.Household, error)
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

func (service *HouseholdService) Search(searchDto dto.SearchHouseholdDto) ([]model.Household, error) {
	households, err := service.repository.Search(searchDto)
	if err != nil {
		return nil, err
	}
	return households, nil
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
		//OwnerID:         1,
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
