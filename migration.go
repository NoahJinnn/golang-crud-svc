package main

import (
	"github.com/ecoprohcm/DMS_BackendServer/models"
	"gorm.io/gorm"
)

func migrateTable(db *gorm.DB, t interface{}) {
	if !db.Migrator().HasTable(t) {
		db.Migrator().CreateTable(t)
	}
}

func migrate(db *gorm.DB) {
	// migrateTable(db, &models.Gateway{})
	// migrateTable(db, &models.Area{})
	// migrateTable(db, &models.Doorlock{})
	// migrateTable(db, &models.GatewayLog{})
	// migrateTable(db, &models.Password{})
	// migrateTable(db, &models.Student{})
	// migrateTable(db, &models.Employee{})
	// migrateTable(db, &models.Scheduler{})
	db.AutoMigrate(&models.Area{}, &models.Gateway{}, &models.Doorlock{}, &models.GatewayLog{}, &models.Password{}, &models.Employee{}, &models.Student{}, &models.Scheduler{})
	// err := db.SetupJoinTable(&models.Student{}, "Schedulers", &models.StudentScheduler{})
	// if err != nil {
	// 	panic(err)
	// }
}
