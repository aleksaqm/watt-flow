package service

import (
	"log"
	"os"

	"watt-flow/db"
)

type RestartService struct {
	database    db.Database
	UserService IUserService
}

func (service RestartService) ResetDatabase() error {
	err := service.database.TruncateAllTables()
	// err := service.database.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE").Error
	if err != nil {
		return err
	}
	log.Println("All tables have been reset.")
	return nil
}

func (service RestartService) InitSuperAdmin() error {
	password, err := service.UserService.CreateSuperAdmin()
	if err != nil {
		return err
	}
	file, err2 := os.Create("/app/data/admin_password.txt")
	if err2 != nil {
		return err2
	}
	_, err3 := file.WriteString(password)
	if err3 != nil {
		return err3
	}
	// insert default pricelist valid from one year in the past
	service.database.DB.Exec("INSERT INTO pricelists (valid_from, blue_zone, red_zone, green_zone, billing_power, tax, is_active) VALUES (CURRENT_DATE - INTERVAL '1 year', 0.15, 0.20, 0.10, 5.0, 10, TRUE) ON CONFLICT DO NOTHING")
	return nil
}

func NewRestartService(database db.Database, userService IUserService) *RestartService {
	return &RestartService{
		database:    database,
		UserService: userService,
	}
}
