package models

import (
	"context"
	"time"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Doorlock struct {
	GormModel
	AreaID          string    `gorm:"type:varchar(256)" json:"areaId"`
	GatewayID       string    `gorm:"type:varchar(256)" json:"gatewayId"`
	SchedulerID     string    `gorm:"type:varchar(256)" json:"schedulerId"`
	Description     string    `json:"description"`
	Location        string    `gorm:"unique;not null" json:"location"`
	LastConnectTime time.Time `json:"lastConnectTime"`
	State           string    `gorm:"not null" json:"state"`
}

type DoorlockCmd struct {
	ID      string `json:"id"`
	Command string `json:"command"`
}

type DoorlockSvc struct {
	db *gorm.DB
}

func NewDoorlockSvc(db *gorm.DB) *DoorlockSvc {
	return &DoorlockSvc{
		db: db,
	}
}

func (dl *Doorlock) BeforeCreate(tx *gorm.DB) (err error) {
	dl.ID = uuid.New().String()
	return
}

func (dls *DoorlockSvc) FindAllDoorlock(ctx context.Context) (dlList []Doorlock, err error) {
	result := dls.db.Find(&dlList)
	if err := result.Error; err != nil {
		err = utils.QueryErrorHandler(err)
		return nil, err
	}
	return dlList, nil
}

func (dls *DoorlockSvc) FindDoorlockByID(ctx context.Context, id string) (dl *Doorlock, err error) {
	result := dls.db.First(&dl, id)
	if err := result.Error; err != nil {
		err = utils.QueryErrorHandler(err)
		return nil, err
	}
	return dl, nil
}

func (dls *DoorlockSvc) CreateDoorlock(ctx context.Context, dl *Doorlock) (*Doorlock, error) {
	if err := dls.db.Create(&dl).Error; err != nil {
		err = utils.QueryErrorHandler(err)
		return nil, err
	}
	return dl, nil
}

func (dls *DoorlockSvc) UpdateDoorlock(ctx context.Context, dl *Doorlock) (bool, error) {
	result := dls.db.Model(&dl).Where("id = ?", dl.ID).Updates(dl)
	return utils.ReturnBoolStateFromResult(result)
}

func (dls *DoorlockSvc) UpdateDoorlockGateway(ctx context.Context, dl *Doorlock, gwID string) (bool, error) {
	result := dls.db.Model(&dl).Where("id = ?", dl.ID).Update("gateway_id", gwID)
	return utils.ReturnBoolStateFromResult(result)
}

func (dls *DoorlockSvc) DeleteDoorlock(ctx context.Context, dl *Doorlock) (bool, error) {
	result := dls.db.Unscoped().Where("id = ?", dl.ID).Delete(dl)
	return utils.ReturnBoolStateFromResult(result)
}
