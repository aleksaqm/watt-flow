package service

import (
	"errors"
	"gorm.io/gorm"
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/repository"
)

type IHouseholdAccessService interface {
	GrantAccess(householdID uint64, userIDToGrant uint64, currentUserID uint64) error
	RevokeAccess(householdID uint64, userIDToRevoke uint64, currentUserID uint64) error
	GetUsersWithAccess(householdID uint64, currentUserID uint64) ([]dto.UserDto, error)
}

type HouseholdAccessService struct {
	householdAccessRepository repository.HouseholdAccessRepository
	householdRepository       repository.HouseholdRepository
	userRepository            repository.UserRepository
}

func NewHouseholdAccessService(accessRepo repository.HouseholdAccessRepository, householdRepo repository.HouseholdRepository, userRepo repository.UserRepository) *HouseholdAccessService {
	return &HouseholdAccessService{
		householdAccessRepository: accessRepo,
		householdRepository:       householdRepo,
		userRepository:            userRepo,
	}
}

func (s *HouseholdAccessService) GrantAccess(householdID uint64, userIDToGrant uint64, currentUserID uint64) error {
	household, err := s.householdRepository.FindById(householdID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("household not found")
		}
		return err
	}

	if *household.OwnerID != currentUserID {
		return errors.New("forbidden: only the owner can grant access")
	}

	if userIDToGrant == *household.OwnerID {
		return errors.New("owner already has full access")
	}

	existingAccess, err := s.householdAccessRepository.FindByHouseholdIDAndUserID(householdID, userIDToGrant)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingAccess != nil {
		return errors.New("user already has access to this household")
	}

	access := &model.HouseholdAccess{
		HouseholdID: householdID,
		UserID:      userIDToGrant,
	}

	return s.householdAccessRepository.Create(access)
}

func (s *HouseholdAccessService) RevokeAccess(householdID uint64, userIDToRevoke uint64, currentUserID uint64) error {
	household, err := s.householdRepository.FindById(householdID)
	if err != nil {
		return errors.New("household not found")
	}
	if *household.OwnerID != currentUserID {
		return errors.New("forbidden: only the owner can revoke access")
	}

	return s.householdAccessRepository.Delete(householdID, userIDToRevoke)
}

func (s *HouseholdAccessService) GetUsersWithAccess(householdID uint64, currentUserID uint64) ([]dto.UserDto, error) {
	household, err := s.householdRepository.FindById(householdID)
	if err != nil {
		return nil, errors.New("household not found")
	}
	if *household.OwnerID != currentUserID {
		return nil, errors.New("forbidden: only the owner can view access list")
	}

	accesses, err := s.householdAccessRepository.FindByHouseholdID(householdID)
	if err != nil {
		return nil, err
	}

	var results []dto.UserDto
	for _, access := range accesses {
		results = append(results, dto.UserDto{
			Id:       access.User.Id,
			Username: access.User.Username,
			Email:    access.User.Email,
			Role:     access.User.Role.RoleToString(),
			Status:   access.User.Status.StatusToString(),
		})
	}
	return results, nil
}
