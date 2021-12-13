package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Gateway struct {
	gorm.Model
	GatewayID string `gorm:"primaryKey" json:"gateway_id"`
	AreaID    uint
	Name      string `gorm:"unique;not null" json:"name"`
}

func (g *Gateway) CreateGateway(db *gorm.DB) error {
	res := db.Create(g)
	fmt.Println(res)
	return nil
}
