package models

import (
	"context"
	"fmt"

	"github.com/trancongduynguyen1997/golang-crud-svc/utils"
	"gorm.io/gorm"
)

type Gateway struct {
	GormModel
	AreaID    uint       `json:"areaId"`
	GatewayID string     `gorm:"type:varchar(256);unique;not null;" json:"gatewayId"`
	Name      string     `json:"name"`
	Doorlocks []Doorlock `gorm:"foreignKey:GatewayID;references:GatewayID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"doorlocks"`
}

// Struct defines HTTP request payload for deleting gateway
type DeleteGateway struct {
	GatewayID string `json:"gatewayId" binding:"required"`
}

type GatewaySvc struct {
	db *gorm.DB
}

func NewGatewaySvc(db *gorm.DB) *GatewaySvc {
	return &GatewaySvc{
		db: db,
	}
}

func (gs *GatewaySvc) FindAllGateway(ctx context.Context) (gwList []Gateway, err error) {
	result := gs.db.Preload("Doorlocks").Find(&gwList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return gwList, nil
}

func (gs *GatewaySvc) FindGatewayByID(ctx context.Context, id string) (gw *Gateway, err error) {
	result := gs.db.Preload("Doorlocks").First(&gw, id)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return gw, nil
}

func (gs *GatewaySvc) FindGatewayByMacID(ctx context.Context, id string) (gw *Gateway, err error) {
	var cnt int64
	result := gs.db.Where("gateway_id = ?", id).Find(&gw).Count(&cnt)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	if cnt <= 0 {
		return nil, fmt.Errorf("find no records")
	}

	return gw, nil
}

func (gs *GatewaySvc) CreateGateway(ctx context.Context, g *Gateway) (*Gateway, error) {
	if err := gs.db.Create(&g).Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return g, nil
}

func (gs *GatewaySvc) UpdateGateway(ctx context.Context, g *Gateway) (bool, error) {
	result := gs.db.Model(&g).Where("id = ?", g.ID).Updates(g)
	return utils.ReturnBoolStateFromResult(result)

}

func (gs *GatewaySvc) DeleteGateway(ctx context.Context, gwID string) (bool, error) {
	result := gs.db.Unscoped().Where("gateway_id = ?", gwID).Delete(&Gateway{})
	return utils.ReturnBoolStateFromResult(result)
}

func (gs *GatewaySvc) AppendGatewayDoorlock(ctx context.Context, gw *Gateway, d *Doorlock) (*Gateway, error) {
	if err := gs.db.Model(&gw).Association("Doorlocks").Append(d); err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return gw, nil
}

func (gs *GatewaySvc) UpdateGatewayDoorlock(ctx context.Context, gw *Gateway, d *Doorlock) (*Gateway, error) {
	if err := gs.db.Model(&gw).Association("Doorlocks").Replace(d); err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return gw, nil
}

func (gs *GatewaySvc) DeleteGatewayDoorlock(ctx context.Context, gw *Gateway, d *Doorlock) (*Gateway, error) {
	if err := gs.db.Model(&gw).Association("Doorlocks").Delete(d); err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return gw, nil
}
