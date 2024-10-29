package repository

import (
	"watt-flow/db"
	"watt-flow/util"
)

type UserRepository struct {
	database db.Database
	logger   util.Logger
}

func NewUserRepository(db db.Database, logger util.Logger) UserRepository {
	return UserRepository{
		database: db,
		logger:   logger,
	}
}
