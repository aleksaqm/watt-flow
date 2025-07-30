package service

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/repository"
	"watt-flow/util"

	"github.com/google/uuid"
)

type IPropertyService interface {
	FindById(id uint64) (*model.Property, error)
	Create(property *model.Property) (*model.Property, error)
	Update(property *model.Property) (*model.Property, error)
	Delete(id uint64) error
	FindByStatus(status model.PropertyStatus) ([]model.Property, error)
	TableQuery(params *dto.PropertyQueryParams) ([]model.Property, int64, error)
	AcceptProperty(id uint64) error
	DeclineProperty(id uint64, message string) error
}

type PropertyService struct {
	propertyRepository repository.PropertyRepository
	householdService   IHouseholdService
	emailSender        *util.EmailSender
}

func NewPropertyService(propertyRepository repository.PropertyRepository, householdService IHouseholdService, emailSender *util.EmailSender) *PropertyService {
	return &PropertyService{
		propertyRepository: propertyRepository,
		householdService:   householdService,
		emailSender:        emailSender,
	}
}

func (service *PropertyService) FindById(id uint64) (*model.Property, error) {
	property, err := service.propertyRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	return property, nil
}

func (service *PropertyService) FindByStatus(status model.PropertyStatus) ([]model.Property, error) {
	properties, err := service.propertyRepository.FindByStatus(status)
	if err != nil {
		return nil, err
	}
	return properties, nil
}

func (service *PropertyService) Create(property *model.Property) (*model.Property, error) {
	UUID := uuid.New()
	var savedFilePaths []string
	prefix := "/app/data/"
	if len(property.Images) > 0 {
		for i, base64String := range property.Images {
			if strings.HasPrefix(base64String, "data:image/") {
				base64String = strings.SplitN(base64String, ",", 2)[1]
			}
			filePath, err := util.SaveFile(UUID.String()+"-"+strconv.Itoa(i), base64String, "jpg", "property_images")
			if err != nil {
				service.cleanupFiles(savedFilePaths)
				return nil, fmt.Errorf("failed to save image %d: %v", i, err)
			}
			property.Images[i] = strings.TrimPrefix(filePath, prefix)
			savedFilePaths = append(savedFilePaths, filePath)
		}
	}

	if len(property.Documents) > 0 {
		for i, base64String := range property.Documents {
			if strings.HasPrefix(base64String, "data:application/") {
				base64String = strings.SplitN(base64String, ",", 2)[1]
			}
			filePath, err := util.SaveFile(UUID.String()+"-"+strconv.Itoa(i), base64String, "pdf", "property_documents")
			if err != nil {
				service.cleanupFiles(savedFilePaths)
				return nil, fmt.Errorf("failed to save document %d: %v", i, err)
			}
			property.Documents[i] = strings.TrimPrefix(filePath, prefix)
			savedFilePaths = append(savedFilePaths, filePath)
		}
	}

	createdProperty, err := service.propertyRepository.Create(property)
	if err != nil {
		service.cleanupFiles(savedFilePaths)
		return nil, fmt.Errorf("failed to create property: %v", err)
	}

	return &createdProperty, nil
}

func (service *PropertyService) cleanupFiles(paths []string) {
	for _, path := range paths {
		err := os.Remove(path)
		if err != nil {
			service.propertyRepository.Logger.Error(err)
		}
	}
}

func (service *PropertyService) Update(property *model.Property) (*model.Property, error) {
	updatedProperty, err := service.propertyRepository.Update(property)
	if err != nil {
		return nil, err
	}
	return &updatedProperty, nil
}

func (service *PropertyService) Delete(id uint64) error {
	err := service.propertyRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (service *PropertyService) TableQuery(params *dto.PropertyQueryParams) ([]model.Property, int64, error) {
	properties, total, err := service.propertyRepository.TableQuery(params)
	if err != nil {
		return nil, 0, err
	}
	return properties, total, nil
}

func (service *PropertyService) AcceptProperty(id uint64) error {
	tx := service.propertyRepository.Database.Begin()
	if tx.Error != nil {
		service.propertyRepository.Logger.Error("Error starting transaction", tx.Error)
		return tx.Error
	}

	property, err := service.FindById(id)
	if err != nil {
		tx.Rollback()
		service.propertyRepository.Logger.Error("Error finding property ", err)
		return err
	}

	err = service.propertyRepository.AcceptProperty(tx, id)
	if err != nil {
		tx.Rollback()
		service.propertyRepository.Logger.Error("Error updating property status", err)
		return err
	}

	err = service.householdService.AcceptHouseholds(tx, id)
	if err != nil {
		tx.Rollback()
		service.propertyRepository.Logger.Error("Error updating households", err)
		return err
	}

	tx.Commit()

	service.propertyRepository.Logger.Info(fmt.Sprintf("Property and its households updated to status for property ID %d", id))

	emailBody := util.GeneratePropertyApprovalEmailBody(property.Address.City+", "+property.Address.Street+" "+property.Address.Number,
		"http://localhost:5173/")

	err = service.emailSender.SendEmail(property.Owner.Email, "Property approved", emailBody)
	return err
}

func (service *PropertyService) DeclineProperty(id uint64, message string) error {
	err := service.propertyRepository.DeclineProperty(id)
	if err != nil {
		service.propertyRepository.Logger.Error("Error updating property status", err)
		return err
	}

	property, err := service.FindById(id)
	if err != nil {
		service.propertyRepository.Logger.Error("Error finding property ", err)
		return err
	}

	emailBody := util.GeneratePropertyDeclineEmailBody(property.Address.City+", "+property.Address.Street+" "+property.Address.Number,
		message, "http://localhost:5173/")

	err = service.emailSender.SendEmail(property.Owner.Email, "Property declined", emailBody)
	return err
}
