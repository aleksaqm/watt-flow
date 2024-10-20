package db

import (
	"fmt"
	"watt-flow/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func NewDatabase() Database {
	c := config.GetConfig()

	url := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", c.GetString("database.host"),
		c.GetString("database.user"), c.GetString("database.password"), c.GetString("database.dbname"), c.GetString("database.port"))

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  url,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	return Database{
		DB: db,
	}
}
