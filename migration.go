package main

import (
	"github.com/ecoprohcm/DMS_BackendServer/models"
	"gorm.io/gorm"
)

func migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.Area{},
		&models.Gateway{},
		&models.Doorlock{},
		&models.GatewayLog{},
		&models.Employee{},
		&models.Student{},
		&models.Scheduler{},
	)
	if err != nil {
		panic(err)
	}
}
