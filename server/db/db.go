package db

import (
	"fmt"
	"log"
	"watt-flow/config"
	"watt-flow/model"
	"watt-flow/util"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func NewDatabase(env *config.Environment, logger util.Logger) Database {
	url := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", env.DBHost,
		env.DBUsername, env.DBPassword, env.DBName, env.DBPort)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  url,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.GetGormLogger(),
	})
	if err != nil {
		fmt.Println(err)
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	return Database{
		DB: db,
	}
}

func (db Database) TruncateAllTables() error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Disable foreign key checks during truncation
		if err := tx.Exec("SET CONSTRAINTS ALL DEFERRED").Error; err != nil {
			return err
		}

		// List of all tables in order considering foreign key dependencies
		tables := []interface{}{
			&model.Household{},
			&model.DeviceStatus{},
			&model.Property{},
			&model.User{},
			&model.Address{},
			&model.City{},
		}

		// Truncate each table
		for _, table := range tables {
			if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(table).Error; err != nil {
				return err
			}
		}

		// Re-enable foreign key checks
		if err := tx.Exec("SET CONSTRAINTS ALL IMMEDIATE").Error; err != nil {
			return err
		}

		return nil
	})
}
