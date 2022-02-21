package main

import (
	"github.com/trancongduynguyen1997/golang-crud-svc/models"
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
		&models.Customer{},
		&models.Scheduler{},
	)
	if err != nil {
		panic(err)
	}
}
