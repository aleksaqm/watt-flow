package service

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"strconv"
	"strings"
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/repository"
	"watt-flow/util"
)

type IPropertyService interface {
	FindById(id uint64) (*model.Property, error)
	Create(property *model.Property) (*model.Property, error)
	Update(property *model.Property) (*model.Property, error)
	Delete(id uint64) error
	FindByStatus(status model.PropertyStatus) ([]model.Property, error)
	TableQuery(params *dto.PropertyQueryParams) ([]model.Property, int64, error)
}

type PropertyService struct {
	propertyRepository *repository.PropertyRepository
	householdService   IHouseholdService
}

func NewPropertyService(propertyRepository *repository.PropertyRepository, householdService IHouseholdService) *PropertyService {
	return &PropertyService{
		propertyRepository: propertyRepository,
		householdService:   householdService,
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
			fullPath := filePath + "/" + UUID.String() + "-" + strconv.Itoa(i) + ".jpg"
			property.Images[i] = fullPath
			savedFilePaths = append(savedFilePaths, fullPath)
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
			fullPath := filePath + "/" + UUID.String() + "-" + strconv.Itoa(i) + ".pdf"
			property.Documents[i] = fullPath
			savedFilePaths = append(savedFilePaths, fullPath)
		}
	}

	createdProperty, err := service.propertyRepository.Create(property)
	if err != nil {
		service.cleanupFiles(savedFilePaths)
		return nil, fmt.Errorf("failed to create property: %v", err)
	}

	return &createdProperty, nil
}

func (service *PropertyService) SearchHouseholds(searchDto dto.SearchHouseholdDto) ([]model.Household, error) {
	var households []model.Household
	if searchDto.Id != "" {
		household, err := service.householdService.FindByCadastralNumber(searchDto.Id)
		if err != nil {
			return nil, err
		}
		households = make([]model.Household, 0)
		households = append(households, *household)
		return households, nil
	}
	households, err := service.householdService.Search(searchDto)
	if err != nil {
		return nil, err
	}
	return households, nil
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
