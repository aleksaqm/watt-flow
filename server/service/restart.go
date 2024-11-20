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
	return nil
}

func NewRestartService(database db.Database, userService IUserService) *RestartService {
	return &RestartService{
		database:    database,
		UserService: userService,
	}
}
