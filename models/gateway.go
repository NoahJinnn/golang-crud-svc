package models

import (
	"context"
	"fmt"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type Gateway struct {
	gorm.Model
	AreaID    uint   `json:"areaId"`
	GatewayID string `gorm:"unique;not null" json:"gatewayId"`
	Name      string `gorm:"unique;not null" json:"name"`
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
	result := gs.db.Find(&gwList)
	if err := result.Error; err != nil {
		err = utils.QueryErrorHandler(err)
		return nil, err
	}
	return gwList, nil
}

func (gs *GatewaySvc) FindGatewayByID(ctx context.Context, id uint) (gw *Gateway, err error) {
	result := gs.db.First(&gw, id)
	if err := result.Error; err != nil {
		err = utils.QueryErrorHandler(err)
		return nil, err
	}
	return gw, nil
}

func (gs *GatewaySvc) CreateGateway(ctx context.Context, g *Gateway) (*Gateway, error) {
	if err := gs.db.Create(&g).Error; err != nil {
		err = utils.QueryErrorHandler(err)
		return nil, err
	}
	return g, nil
}

func (gs *GatewaySvc) UpdateGateway(ctx context.Context, g *Gateway) (*Gateway, error) {
	result := gs.db.Model(&g).Where("id = ?", g.ID).Updates(g)
	err := result.Error
	ra := result.RowsAffected
	if err != nil {
		err = utils.QueryErrorHandler(err)
		return nil, err
	}
	if ra > 0 {
		return g, nil
	} else {
		return nil, fmt.Errorf("no record affected")
	}
}

func (gs *GatewaySvc) DeleteGateway(ctx context.Context, g *Gateway) (bool, error) {
	result := gs.db.Unscoped().Where("id = ?", g.ID).Delete(g)
	return utils.DeleteResultHandler(result)
}
