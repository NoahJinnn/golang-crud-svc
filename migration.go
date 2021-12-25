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
	migrateTable(db, &models.Gateway{})
	migrateTable(db, &models.Area{})
	migrateTable(db, &models.Doorlock{})
	migrateTable(db, &models.GatewayLog{})
}
