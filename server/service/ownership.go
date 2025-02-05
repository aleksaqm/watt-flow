package service

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/repository"
	"watt-flow/util"

	"github.com/google/uuid"
)

type IOwnershipService interface {
	CreateOwnershipRequest(ownershipRequest dto.OwnershipRequestDto) (*dto.OwnershipRequestDto, error)
	GetOwnersRequests(ownerId uint64, params *dto.OwnershipQueryParams) ([]dto.OwnershipResponseDto, int64, error)
	AcceptOwnershipRequest(id uint64) error
	DeclineOwnershipRequest(id uint64, reason string) error
	GetPendingRequests(params *dto.OwnershipQueryParams) ([]dto.OwnershipResponseDto, int64, error)
}

type OwnershipService struct {
	householdRepository repository.HouseholdRepository
	ownershipRepository repository.OwnershipRepository
	emailSender         *util.EmailSender
}

func NewOwnershipService(householdRepository repository.HouseholdRepository, ownershipRepository repository.OwnershipRepository, emailSender *util.EmailSender) *OwnershipService {
	return &OwnershipService{
		householdRepository: householdRepository,
		ownershipRepository: ownershipRepository,
		emailSender:         emailSender,
	}
}

func (service *OwnershipService) CreateOwnershipRequest(ownershipRequestDto dto.OwnershipRequestDto) (*dto.OwnershipRequestDto, error) {
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
	prefix := "/app/data/"
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
			ownershipRequest.Images[i] = strings.TrimPrefix(filePath, prefix)
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
			ownershipRequest.Documents[i] = strings.TrimPrefix(filePath, prefix)
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

func (service *OwnershipService) GetOwnersRequests(ownerId uint64, params *dto.OwnershipQueryParams) ([]dto.OwnershipResponseDto, int64, error) {
	requests, total, err := service.ownershipRepository.FindForOwner(ownerId, params)
	if err != nil {
		return nil, 0, err
	}
	results := make([]dto.OwnershipResponseDto, 0)
	for _, result := range requests {
		mappedRequest, _ := service.MapToOwnershipDto(&result)
		results = append(results, mappedRequest)
	}
	return results, total, nil
}

func (service *OwnershipService) GetPendingRequests(params *dto.OwnershipQueryParams) ([]dto.OwnershipResponseDto, int64, error) {
	requests, total, err := service.ownershipRepository.FindPendingRequests(params)
	if err != nil {
		return nil, 0, err
	}
	results := make([]dto.OwnershipResponseDto, 0)
	for _, result := range requests {
		mappedRequest, _ := service.MapToOwnershipDto(&result)
		results = append(results, mappedRequest)
	}
	return results, total, nil
}

func (service *OwnershipService) AcceptOwnershipRequest(id uint64) error {
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

	if request.Status != 0 {
		service.ownershipRepository.Logger.Error("Error request isn't pending", tx.Error)
		tx.Rollback()
		return errors.New("request isn't pending")
	}

	err = service.ownershipRepository.AcceptRequest(tx, id)
	if err != nil {
		tx.Rollback()
		service.ownershipRepository.Logger.Error("Error accepting request", tx.Error)
		return err
	}

	emailForDenial, err := service.ownershipRepository.DeclineAllForHousehold(tx, request.HouseholdID)
	if err != nil {
		tx.Rollback()
		service.ownershipRepository.Logger.Error("Error accepting request", tx.Error)
		return err
	}

	err = service.householdRepository.AddOwnerToHousehold(tx, request.HouseholdID, request.OwnerID)
	if err != nil {
		tx.Rollback()
		service.householdRepository.Logger.Error("Error adding owner to household", tx.Error)
		return err
	}
	tx.Commit()
	emailBody := util.GenerateOwnershipApprovalEmailBody(request.Household.Property.Address.City+", "+request.Household.Property.Address.Street+" "+request.Household.Property.Address.Number+" suite: "+request.Household.Suite,
		"http://localhost:5173/")

	err = service.emailSender.SendEmail(request.Owner.Email, "Ownership approved", emailBody)

	for _, s := range emailForDenial {
		emailBody := util.GenerateOwnershipDenialEmailBody(request.Household.Property.Address.City+", "+request.Household.Property.Address.Street+" "+request.Household.Property.Address.Number+" suite: "+request.Household.Suite, "We accepted someone else's request.",
			"http://localhost:5173/")
		err = service.emailSender.SendEmail(s, "Ownership declined", emailBody)
	}

	return nil
}

func (service *OwnershipService) DeclineOwnershipRequest(id uint64, reason string) error {
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
	if request.Status != 0 {
		service.ownershipRepository.Logger.Error("Error request isn't pending", tx.Error)
		tx.Rollback()
		return errors.New("request isn't pending")
	}

	err = service.ownershipRepository.DeclineRequest(tx, id, reason)
	if err != nil {
		tx.Rollback()
		service.ownershipRepository.Logger.Error("Error accepting request", tx.Error)
		return err
	}
	tx.Commit()
	emailBody := util.GenerateOwnershipDenialEmailBody(request.Household.Property.Address.City+", "+request.Household.Property.Address.Street+" "+request.Household.Property.Address.Number+" suite: "+request.Household.Suite, reason,
		"http://localhost:5173/")
	err = service.emailSender.SendEmail(request.Owner.Email, "Ownership declined", emailBody)

	return nil
}

func (service *OwnershipService) cleanupFiles(paths []string) {
	for _, path := range paths {
		err := os.Remove(path)
		if err != nil {
			service.ownershipRepository.Logger.Error(err)
		}
	}
}

func (service *OwnershipService) MapToOwnershipDto(ownership *model.OwnershipRequest) (dto.OwnershipResponseDto, error) {
	response := dto.OwnershipResponseDto{
		Id:          ownership.Id,
		Images:      ownership.Images,
		Documents:   ownership.Documents,
		OwnerID:     ownership.OwnerID,
		HouseholdID: ownership.HouseholdID,
		CreatedAt:   ownership.CreatedAt.String(),
		ClosedAt:    ownership.ClosedAt.String(), // mozda pukne null pointer exception
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
