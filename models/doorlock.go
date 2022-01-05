package models

import (
	"context"
	"time"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type Doorlock struct {
	GormModel
	DoorSerialID    string      `gorm:"unique;not null" json:"doorSerialId"`
	Location        string      `gorm:"unique;not null" json:"location"`
	State           string      `gorm:"not null" json:"state"`
	Description     string      `json:"description"`
	LastConnectTime time.Time   `json:"lastConnectTime"`
	AreaID          uint        `json:"areaId"`
	GatewayID       string      `gorm:"type:varchar(256);" json:"gatewayId"`
	Schedulers      []Scheduler `gorm:"many2many:door_schedulers;"`
}

type DoorlockCmd struct {
	DoorSerialID string `json:"doorSerialId"`
	GatewayID    string `json:"gatewayId"`
	Command      string `json:"command"`
}

type DoorlockDelete struct {
	DoorSerialID string `json:"doorSerialId"`
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

func (d *Doorlock) BeforeCreate(tx *gorm.DB) (err error) {
	if d.State == "" {
		d.State = "Close"
	}
	return
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

func (dls *DoorlockSvc) UpdateDoorlockGateway(ctx context.Context, dl *Doorlock, gwID string) (bool, error) {
	result := dls.db.Model(&dl).Where("id = ?", dl.ID).Update("gateway_id", gwID)
	return utils.ReturnBoolStateFromResult(result)
}

func (dls *DoorlockSvc) DeleteDoorlock(ctx context.Context, doorSerialId string) (bool, error) {
	result := dls.db.Unscoped().Where("door_serial_id = ?", doorSerialId).Delete(&Doorlock{})
	return utils.ReturnBoolStateFromResult(result)
}
