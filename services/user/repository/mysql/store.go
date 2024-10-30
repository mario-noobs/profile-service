package mysql

import (
	"demo-service/helpers"
	"gorm.io/gorm"
)

type mysqlRepo struct {
	db   *gorm.DB
	time helpers.Timer
}

func NewMySQLRepository(db *gorm.DB) *mysqlRepo {
	return &mysqlRepo{db: db}
}
