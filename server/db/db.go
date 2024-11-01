package db

import (
	"fmt"

	"watt-flow/config"
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
	return Database{
		DB: db,
	}
}
