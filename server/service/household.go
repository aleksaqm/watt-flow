package service

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/repository"
	"watt-flow/util"

	"gorm.io/gorm"
)

type IHouseholdService interface {
	FindById(id uint64) (*dto.HouseholdResultDto, error)
	Create(household *dto.CreateHouseholdDto) (*model.Household, error)
	Update(household *model.Household) (*model.Household, error)
	Delete(id uint64) error
	FindByStatus(status model.HouseholdStatus) ([]model.Household, error)
	FindByCadastralNumber(id string) (*model.Household, error)
	Query(queryParams *dto.HouseholdQueryParams) ([]dto.HouseholdResultDto, int64, error)
	AcceptHouseholds(tx *gorm.DB, propertyID uint64) error
	CreateOwnershipRequest(ownershipRequest dto.OwnershipRequestDto) (*dto.OwnershipRequestDto, error)
	GetOwnersRequests(ownerId uint64, params *dto.OwnershipQueryParams) ([]dto.OwnershipResponseDto, int64, error)
	AcceptOwnershipRequest(id uint64) error
}

type HouseholdService struct {
	repository          *repository.HouseholdRepository
	ownershipRepository *repository.OwnershipRepository
}

func NewHouseholdService(repository *repository.HouseholdRepository, ownershipRepository *repository.OwnershipRepository) *HouseholdService {
	return &HouseholdService{
		repository:          repository,
		ownershipRepository: ownershipRepository,
	}
}

func (service *HouseholdService) FindById(id uint64) (*dto.HouseholdResultDto, error) {
	household, err := service.repository.FindById(id)
	if err != nil {
		return nil, err
	}
	mappedHousehold, _ := MapToResultDto(household)

	return &mappedHousehold, nil
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
		Latitude:        household.Property.Address.Latitude,
		Longitude:       household.Property.Address.Longitude,
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

func (service *HouseholdService) CreateOwnershipRequest(ownershipRequestDto dto.OwnershipRequestDto) (*dto.OwnershipRequestDto, error) {
	UUID := uuid.New()
	var savedFilePaths []string
	ownershipRequest := model.OwnershipRequest{
		Status:      model.Pending,
		OwnerID:     ownershipRequestDto.OwnerID,
		HouseholdID: ownershipRequestDto.HouseholdID,
		Images:      make([]string, len(ownershipRequestDto.Images)),
		Documents:   make([]string, len(ownershipRequestDto.Documents)),
		CreatedAt:   time.Now(),
	}
	if len(ownershipRequestDto.Images) > 0 {
		for i, base64String := range ownershipRequestDto.Images {
			if strings.HasPrefix(base64String, "data:image/") {
				base64String = strings.SplitN(base64String, ",", 2)[1]
			}
			filePath, err := util.SaveFile(UUID.String()+"-"+strconv.Itoa(i), base64String, "jpg", "ownership_images")
			if err != nil {
				service.cleanupFiles(savedFilePaths)
				return nil, fmt.Errorf("failed to save image %d: %v", i, err)
			}
			//fullPath := filePath + "/" + UUID.String() + "-" + strconv.Itoa(i) + ".jpg"
			ownershipRequest.Images[i] = filePath
			savedFilePaths = append(savedFilePaths, filePath)
		}
	}
	if len(ownershipRequestDto.Documents) > 0 {
		for i, base64String := range ownershipRequestDto.Documents {
			if strings.HasPrefix(base64String, "data:application/") {
				base64String = strings.SplitN(base64String, ",", 2)[1]
			}
			filePath, err := util.SaveFile(UUID.String()+"-"+strconv.Itoa(i), base64String, "pdf", "ownership_documents")
			if err != nil {
				service.cleanupFiles(savedFilePaths)
				return nil, fmt.Errorf("failed to save document %d: %v", i, err)
			}
			//fullPath := filePath + "/" + UUID.String() + "-" + strconv.Itoa(i) + ".pdf"
			ownershipRequest.Documents[i] = filePath
			savedFilePaths = append(savedFilePaths, filePath)
		}
	}
	createdRequest, err := service.ownershipRepository.Create(&ownershipRequest)
	if err != nil {
		return nil, err
	}
	requestDto := dto.OwnershipRequestDto{
		Id:          createdRequest.Id,
		Images:      createdRequest.Images,
		Documents:   createdRequest.Documents,
		OwnerID:     createdRequest.OwnerID,
		HouseholdID: createdRequest.HouseholdID,
		CreatedAt:   createdRequest.CreatedAt.String(),
	}
	return &requestDto, nil
}

func (service *HouseholdService) GetOwnersRequests(ownerId uint64, params *dto.OwnershipQueryParams) ([]dto.OwnershipResponseDto, int64, error) {
	requests, total, err := service.ownershipRepository.FindForOwner(ownerId, params)
	if err != nil {
		return nil, 0, err
	}
	var results = make([]dto.OwnershipResponseDto, 0)
	for _, result := range requests {
		mappedRequest, _ := service.MapToOwnershipDto(&result)
		results = append(results, mappedRequest)
	}
	return results, total, nil
}

func (service *HouseholdService) AcceptOwnershipRequest(id uint64) error {
	tx := service.ownershipRepository.Database.Begin()
	if tx.Error != nil {
		service.ownershipRepository.Logger.Error("Error starting transaction", tx.Error)
		return tx.Error
	}
	request, err := service.ownershipRepository.FindById(id)
	if err != nil {
		tx.Rollback()
		service.ownershipRepository.Logger.Error("Error finding request", tx.Error)
		return err
	}

	err = service.ownershipRepository.AcceptRequest(tx, id)
	if err != nil {
		tx.Rollback()
		service.ownershipRepository.Logger.Error("Error accepting request", tx.Error)
		return err
	}

	err = service.repository.AddOwnerToHousehold(tx, request.HouseholdID, request.OwnerID)
	if err != nil {
		tx.Rollback()
		service.repository.Logger.Error("Error adding owner to household", tx.Error)
		return err
	}
	tx.Commit()

	//salji mejl

	return nil
}

func (service *HouseholdService) cleanupFiles(paths []string) {
	for _, path := range paths {
		err := os.Remove(path)
		if err != nil {
			service.ownershipRepository.Logger.Error(err)
		}
	}
}

func (service *HouseholdService) MapToOwnershipDto(ownership *model.OwnershipRequest) (dto.OwnershipResponseDto, error) {
	service.repository.Logger.Info(ownership.Household)
	service.repository.Logger.Info(ownership.Owner)
	response := dto.OwnershipResponseDto{
		Id:          ownership.Id,
		Images:      ownership.Images,
		Documents:   ownership.Documents,
		OwnerID:     ownership.OwnerID,
		HouseholdID: ownership.HouseholdID,
		CreatedAt:   ownership.CreatedAt.String(),
		ClosedAt:    ownership.ClosedAt.String(), //mozda pukne null pointer exception
		City:        ownership.Household.Property.Address.City,
		Street:      ownership.Household.Property.Address.Street,
		Number:      ownership.Household.Property.Address.Number,
		Floor:       ownership.Household.Floor,
		Suite:       ownership.Household.Suite,
		Username:    ownership.Owner.Username,
		Status:      ownership.Status.RequestStatusToString(),
	}
	return response, nil
}
