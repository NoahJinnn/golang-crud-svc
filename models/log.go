package models

import (
	"context"
	"time"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type GatewayLog struct {
	ID        uint      `gorm:"primarykey"`
	MacID     string    `gorm:"unique; not null" json:"macId"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	CreatedAt time.Time `swaggerignore:"true"`
}

type LogSvc struct {
	db *gorm.DB
}

func NewLogSvc(db *gorm.DB) *LogSvc {
	return &LogSvc{
		db: db,
	}
}

func (gs *GatewaySvc) FindAllGatewayLog(ctx context.Context) (gwList []Gateway, err error) {
	result := gs.db.Preload("Doorlocks").Find(&gwList)
	if err := result.Error; err != nil {
		err = utils.QueryErrorHandler(err)
		return nil, err
	}
	return gwList, nil
}

func (gs *GatewaySvc) FindGatewayLogByID(ctx context.Context, id uint) (gw *Gateway, err error) {
	result := gs.db.Preload("Doorlocks").First(&gw, id)
	if err := result.Error; err != nil {
		err = utils.QueryErrorHandler(err)
		return nil, err
	}
	return gw, nil
}
