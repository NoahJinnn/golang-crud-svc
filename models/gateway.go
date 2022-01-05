package models

import (
	"context"
	"time"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type Gateway struct {
	GormModel
	AreaID          uint       `json:"areaId"`
	GatewayID       string     `gorm:"type:varchar(256);unique;not null;" json:"gatewayId"`
	Name            string     `gorm:"unique;not null" json:"name"`
	LastConnectTime time.Time  `json:"lastConnectTime"`
	State           string     `gorm:"not null" json:"state"`
	Doorlocks       []Doorlock `gorm:"foreignKey:GatewayID;references:GatewayID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"doorlocks"`
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

func (gs *GatewaySvc) DeleteGateway(ctx context.Context, gwId uint) (bool, error) {
	result := gs.db.Unscoped().Where("id = ?", gwId).Delete(&Gateway{})
	return utils.ReturnBoolStateFromResult(result)
}

func (gs *GatewaySvc) DeleteGatewayByMacId(ctx context.Context, g *Gateway) (bool, error) {
	result := gs.db.Unscoped().Where("gateway_id = ?", g.GatewayID).Delete(g)
	return utils.ReturnBoolStateFromResult(result)
}
