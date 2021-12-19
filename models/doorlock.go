package models

import (
	"context"
	"time"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type Doorlock struct {
	gorm.Model
	AreaID          uint      `json:"areaId"`
	GatewayID       uint      `json:"gatewayId"`
	SchedulerID     uint      `json:"schedulerId"`
	Description     string    `json:"description"`
	Location        string    `gorm:"unique;not null" json:"location"`
	LastConnectTime time.Time `json:"lastConnectTime"`
	State           string    `gorm:"not null" json:"state"`
}

type DoorlockSvc struct {
	db *gorm.DB
}

func NewDoorlockSvc(db *gorm.DB) *DoorlockSvc {
	return &DoorlockSvc{
		db: db,
	}
}

func (dls *DoorlockSvc) FindAllDoorlock(ctx context.Context) (dlList []Doorlock, err error) {
	result := dls.db.Find(&dlList)
	if err := result.Error; err != nil {
		err = utils.QueryErrorHandler(err)
		return nil, err
	}
	return dlList, nil
}

func (dls *DoorlockSvc) FindDoorlockByID(ctx context.Context, id uint) (dl *Doorlock, err error) {
	result := dls.db.First(&dl, id)
	if err := result.Error; err != nil {
		err = utils.QueryErrorHandler(err)
		return nil, err
	}
	return dl, nil
}

func (dls *DoorlockSvc) CreateDoorlock(dl *Doorlock, ctx context.Context) (*Doorlock, error) {
	if err := dls.db.Create(&dl).Error; err != nil {
		err = utils.QueryErrorHandler(err)
		return nil, err
	}
	return dl, nil
}

func (dls *DoorlockSvc) UpdateDoorlock(ctx context.Context, dl *Doorlock) (*Doorlock, error) {
	result := dls.db.Model(&dl).Where("id = ?", dl.ID).Updates(dl)
	handled, err := utils.UpdateResultHandler(result, dl)
	if err != nil {
		return nil, err
	}
	dl = handled.(*Doorlock)
	return dl, nil
}

func (dls *DoorlockSvc) DeleteDoorlock(ctx context.Context, dl *Doorlock) (bool, error) {
	result := dls.db.Unscoped().Where("id = ?", dl.ID).Delete(dl)
	return utils.DeleteResultHandler(result)
}
