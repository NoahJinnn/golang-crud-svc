package models

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Gateway struct {
	gorm.Model
	GatewayID string `gorm:"unique;not null" binding:"required" json:"gatewayId"`
	AreaID    uint   `json:"areaId"`
	Name      string `gorm:"unique;not null" binding:"required" json:"name"`
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("can't find any record")
		}
		return nil, err
	}
	return gwList, nil
}

func (gs *GatewaySvc) CreateGateway(ctx context.Context, g *Gateway) (*Gateway, error) {
	if err := gs.db.Create(&g).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("can't find any record")
		}
		return nil, err
	}
	return g, nil
}

func (gs *GatewaySvc) UpdateGateway(ctx context.Context, g *Gateway) (*Gateway, error) {
	result := gs.db.Model(&g).Where("id = ?", g.ID).Updates(g)
	err := result.Error
	ra := result.RowsAffected
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("can't find any record")
		}
		return nil, err
	}
	if ra > 0 {
		return g, nil
	} else {
		return nil, fmt.Errorf("no record affected")
	}
}

func (gs *GatewaySvc) DeleteGateway(ctx context.Context, dg *DeleteGateway) (bool, error) {
	g := convertToDeleteGw(dg)
	fmt.Println("Delete success", g.GatewayID)
	if err := gs.db.Unscoped().Where("gateway_id = ?", g.GatewayID).Delete(g).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, fmt.Errorf("can't find any record")
		}
		return false, err
	}
	return true, nil
}
