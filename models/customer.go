package models

import (
	"context"
	"fmt"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type Customer struct {
	GormModel
	CCCD  string `gorm:"type:varchar(256); unique; not null;" json:"cccd"  binding:"required"`
	Name  string `json:"name"`
	Phone string `gorm:"type:varchar(50)" json:"phone"`
	UserPass
	Schedulers []Scheduler `gorm:"foreignKey:CustomerID;references:CCCD;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"schedulers"`
}

type CustomerSvc struct {
	db *gorm.DB
}

func NewCustomerSvc(db *gorm.DB) *CustomerSvc {
	return &CustomerSvc{
		db: db,
	}
}

func (cs *CustomerSvc) FindAllCustomer(ctx context.Context) (cList []Customer, err error) {
	result := cs.db.Preload("Schedulers").Find(&cList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return cList, nil
}

func (cs *CustomerSvc) FindCustomerByCCCD(ctx context.Context, cccd string) (c *Customer, err error) {
	var cnt int64
	result := cs.db.Preload("Schedulers").Where("mssv = ?", cccd).Find(&c).Count(&cnt)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	if cnt <= 0 {
		return nil, fmt.Errorf("find no records")
	}

	return c, nil
}

func (cs *CustomerSvc) CreateCustomer(ctx context.Context, c *Customer) (*Customer, error) {
	if err := cs.db.Create(&c).Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return c, nil
}

func (cs *CustomerSvc) UpdateCustomer(ctx context.Context, c *Customer) (bool, error) {
	result := cs.db.Model(&c).Where("id = ?", c.ID).Updates(c)
	return utils.ReturnBoolStateFromResult(result)
}

func (cs *CustomerSvc) DeleteCustomer(ctx context.Context, cId uint) (bool, error) {
	result := cs.db.Unscoped().Where("id = ?", cId).Delete(&Customer{})
	return utils.ReturnBoolStateFromResult(result)
}

func (cs *CustomerSvc) AppendCustomerScheduler(ctx context.Context, c *Customer, doorSerialId string, sche *Scheduler) (*Customer, error) {

	// Add scheduler for door
	var door = &Doorlock{}
	doorResult := cs.db.Where("door_serial_id = ?", doorSerialId).First(door)
	if err := doorResult.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	sche.StudentID = &c.CCCD
	sche.DoorSerialID = &doorSerialId
	if err := cs.db.Create(&sche).Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	if err := cs.db.Model(door).Association("Schedulers").Append(sche); err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	// Add scheduler for customer
	if err := cs.db.Model(&c).Association("Schedulers").Append(sche); err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return c, nil
}

// func (cs *CustomerSvc) UpdateCustomerScheduler(ctx context.Context, c *Customer, doorSerialId string, sche *Scheduler) (*Customer, error) {
// 	// Update scheduler for door
// 	var door = &Doorlock{}
// 	doorResult := cs.db.Where("door_serial_id = ?", doorSerialId).First(door)
// 	if err := doorResult.Error; err != nil {
// 		err = utils.HandleQueryError(err)
// 		return nil, err
// 	}

// 	if err := cs.db.Model(door).Association("Schedulers").Replace(sche); err != nil {
// 		err = utils.HandleQueryError(err)
// 		return nil, err
// 	}

// 	// Update scheduler for customer
// 	if err := cs.db.Model(&c).Association("Schedulers").Replace(sche); err != nil {
// 		err = utils.HandleQueryError(err)
// 		return nil, err
// 	}
// 	return c, nil
// }

// func (cs *CustomerSvc) DeleteCustomerScheduler(ctx context.Context, c *Customer, doorSerialId string, sche *Scheduler) (*Customer, error) {
// 	// Delete scheduler for door
// 	var door = &Doorlock{}
// 	doorResult := cs.db.Where("door_serial_id = ?", doorSerialId).First(door)
// 	if err := doorResult.Error; err != nil {
// 		err = utils.HandleQueryError(err)
// 		return nil, err
// 	}

// 	if err := cs.db.Model(door).Association("Schedulers").Delete(sche); err != nil {
// 		err = utils.HandleQueryError(err)
// 		return nil, err
// 	}

// 	// Delete scheduler for customer
// 	if err := cs.db.Model(&c).Association("Schedulers").Delete(sche); err != nil {
// 		err = utils.HandleQueryError(err)
// 		return nil, err
// 	}
// 	return c, nil
// }
