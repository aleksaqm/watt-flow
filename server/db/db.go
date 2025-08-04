package db

import (
	"fmt"
	"log"
	"time"

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

		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: false,
	})
	if err != nil {
		fmt.Println(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}
	sqlDB.SetMaxOpenConns(300)
	sqlDB.SetMaxIdleConns(30)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	return Database{
		DB: db,
	}
}

func (db Database) TruncateAllTables() error {
	truncateQuery := `
        TRUNCATE TABLE
            public.household_accesses,
            public.bills,
            public.ownership_requests,
            public.meetings,
            public.time_slots,
            public.device_status,
            public.households,
            public.properties,
            public.users,
            public.pricelists,
            public.monthly_bills,
            public.cities,
            public.addresses
        RESTART IDENTITY CASCADE;
    `
	return db.Exec(truncateQuery).Error
}
