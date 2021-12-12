package models

import (
	"gorm.io/gorm"
)

type Gateway struct {
	gorm.Model
	GatewayID string `gorm:"primaryKey" json:"gateway_id"`
	AreaID    uint
	Name      string `gorm:"unique;not null" json:"name"`
}
