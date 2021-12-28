package models

import (
	"context"
	"time"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GatewayLog struct {
	ID        string    `gorm:"primarykey; type:varchar(256)"`
	MacID     string    `gorm:"unique; not null" json:"macId"`
	LogType   string    `json:"logType"`
	Content   string    `json:"content"`
	LogTime   string    `json:"logTime"`
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

func (l *GatewayLog) BeforeCreate(tx *gorm.DB) (err error) {
	l.ID = uuid.New().String()
	return
}

func (ls *LogSvc) FindAllGatewayLog(ctx context.Context) (glList []GatewayLog, err error) {
	result := ls.db.Find(&glList)
	if err := result.Error; err != nil {
		err = utils.QueryErrorHandler(err)
		return nil, err
	}
	return glList, nil
}

func (ls *LogSvc) FindGatewayLogByID(ctx context.Context, id string) (gl *GatewayLog, err error) {
	result := ls.db.Preload("Doorlocks").First(&gl, id)
	if err := result.Error; err != nil {
		err = utils.QueryErrorHandler(err)
		return nil, err
	}
	return gl, nil
}

func (ls *LogSvc) CreateGatewayLog(ctx context.Context, gl *GatewayLog) (*GatewayLog, error) {
	if err := ls.db.Create(&gl).Error; err != nil {
		err = utils.QueryErrorHandler(err)
		return nil, err
	}
	return gl, nil
}
