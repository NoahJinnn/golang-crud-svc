package models

import (
	"context"
	"fmt"
	"time"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type Doorlock struct {
	GormModel
	DoorSerialID string      `gorm:"type:varchar(256);unique;not null" json:"doorSerialId"`
	Location     string      `json:"location"`
	Description  string      `json:"description"`
	GatewayID    string      `gorm:"type:varchar(256);" json:"gatewayId"`
	LastOpenTime uint        `json:"lastOpenTime"`
	ConnectState string      `json:"connectState"`
	Schedulers   []Scheduler `gorm:"foreignKey:DoorSerialID;references:DoorSerialID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"schedulers"`
}

type DoorlockCmd struct {
	DoorSerialID string `json:"doorSerialId"`
	GatewayID    string `json:"gatewayId"`
	State        string `json:"state"`
}

type DoorlockDelete struct {
	DoorSerialID string `json:"doorSerialId" binding:"required"`
	GatewayID    string `json:"gatewayId"`
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
	result := dls.db.Preload("Schedulers").Find(&dlList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dlList, nil
}

func (dls *DoorlockSvc) FindDoorlockByID(ctx context.Context, id string) (dl *Doorlock, err error) {
	result := dls.db.Preload("Schedulers").First(&dl, id)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dl, nil
}

func (dls *DoorlockSvc) FindDoorlockBySerialID(ctx context.Context, serialiId string) (dl *Doorlock, err error) {
	var cnt int64
	result := dls.db.Preload("Schedulers").Where("door_serial_id = ?", serialiId).Find(&dl).Count(&cnt)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	if cnt <= 0 {
		return nil, fmt.Errorf("find no records")
	}

	return dl, nil
}

func (dls *DoorlockSvc) CreateDoorlock(ctx context.Context, dl *Doorlock) (*Doorlock, error) {
	if err := dls.db.Create(&dl).Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return dl, nil
}

func (dls *DoorlockSvc) UpdateDoorlock(ctx context.Context, dl *Doorlock) (bool, error) {
	result := dls.db.Model(&dl).Where("id = ?", dl.ID).Updates(dl)
	return utils.ReturnBoolStateFromResult(result)
}

func (dls *DoorlockSvc) UpdateDoorlockBySerialID(ctx context.Context, dl *Doorlock) (bool, error) {
	result := dls.db.Model(&dl).Where("door_serial_id = ?", dl.DoorSerialID).Updates(dl)
	return utils.ReturnBoolStateFromResult(result)
}

func (dls *DoorlockSvc) UpdateDoorlockState(ctx context.Context, dl *DoorlockCmd) (bool, error) {
	result := dls.db.Model(&Doorlock{}).Where("door_serial_id = ?", dl.DoorSerialID).Update("last_open_time", time.Now().UnixMilli())
	return utils.ReturnBoolStateFromResult(result)
}

func (dls *DoorlockSvc) DeleteDoorlock(ctx context.Context, doorSerialId string) (bool, error) {
	result := dls.db.Unscoped().Where("door_serial_id = ?", doorSerialId).Delete(&Doorlock{})
	return utils.ReturnBoolStateFromResult(result)
}
