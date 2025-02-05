package repository

import (
	"watt-flow/db"
	"watt-flow/model"
	"watt-flow/util"
)

type CityRepository struct {
	database db.Database
	logger   util.Logger
}

func NewCityRepository(db db.Database, logger util.Logger) CityRepository {
	err := db.AutoMigrate(&model.City{})
	if err != nil {
		logger.Error("Error migrating city table", err)
	}
	return CityRepository{
		database: db,
		logger:   logger,
	}
}

func (repository *CityRepository) GetAllCities() ([]string, error) {
	var cities []model.City
	if err := repository.database.Find(&cities).Error; err != nil {
		repository.logger.Error("Error fetching cities", err)
		return nil, err
	}

	var cityNames []string
	for _, city := range cities {
		cityNames = append(cityNames, city.CityName)
	}
	return cityNames, nil
}
