package models

import (
	"context"
	"fmt"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type Employee struct {
	GormModel
	MSNV       string `gorm:"type:varchar(256); unique; not null;" json:"msnv" binding:"required"`
	Name       string `json:"name"`
	Phone      string `gorm:"type:varchar(50)" json:"phone"`
	Email      string `gorm:"type:varchar(256); not null;" json:"email"`
	Department string `json:"department"`
	Role       string `gorm:"not null;" json:"role"`
	UserPass
	HighestPriority bool        `json:"highestPriority"`
	Schedulers      []Scheduler `gorm:"foreignKey:EmployeeID;references:MSNV;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"schedulers"`
}

type DeleteEmployee struct {
	MSNV string `json:"msnv" binding:"required"`
}

type EmployeeSvc struct {
	db *gorm.DB
}

func NewEmployeeSvc(db *gorm.DB) *EmployeeSvc {
	return &EmployeeSvc{
		db: db,
	}
}

func (es *EmployeeSvc) FindAllEmployee(ctx context.Context) (eList []Employee, err error) {
	result := es.db.Preload("Schedulers").Find(&eList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return eList, nil
}

func (es *EmployeeSvc) FindEmployeeByMSNV(ctx context.Context, msnv string) (e *Employee, err error) {
	var cnt int64
	result := es.db.Preload("Schedulers").Where("msnv = ?", msnv).Find(&e).Count(&cnt)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	if cnt <= 0 {
		return nil, fmt.Errorf("find no records")
	}

	return e, nil
}

func (es *EmployeeSvc) FindAllHPEmployee(ctx context.Context) (eL []Employee, err error) {
	result := es.db.Where("highest_priority = ?", true).Find(&eL)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return eL, nil
}

func (es *EmployeeSvc) CreateEmployee(ctx context.Context, e *Employee) (*Employee, error) {
	if err := es.db.Create(&e).Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return e, nil
}

func (es *EmployeeSvc) UpdateEmployee(ctx context.Context, e *Employee) (bool, error) {
	result := es.db.Model(&e).Where("id = ?", e.ID).Updates(e)
	return utils.ReturnBoolStateFromResult(result)
}

func (es *EmployeeSvc) UpdateHPEmployee(ctx context.Context, e *Employee) (bool, error) {
	result := es.db.Model(&e).Where("id = ? AND highest_priority = ?", e.ID, true).Updates(e)
	return utils.ReturnBoolStateFromResult(result)
}

func (es *EmployeeSvc) DeleteEmployee(ctx context.Context, msnv string) (bool, error) {
	result := es.db.Unscoped().Where("msnv = ?", msnv).Delete(&Employee{})
	return utils.ReturnBoolStateFromResult(result)
}

func (es *EmployeeSvc) DeleteHPEmployee(ctx context.Context, msnv string) (bool, error) {
	result := es.db.Unscoped().Where("msnv = ? AND highest_priority = ?", msnv, true).Delete(&Employee{})
	return utils.ReturnBoolStateFromResult(result)
}

func (es *EmployeeSvc) AppendEmployeeScheduler(ctx context.Context, e *Employee, doorSerialId string, sche *Scheduler) (*Employee, error) {

	// Add scheduler for door
	var door = &Doorlock{}
	doorResult := es.db.Where("door_serial_id = ?", doorSerialId).First(door)
	if err := doorResult.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	if err := es.db.Model(door).Association("Schedulers").Append(sche); err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	// Add scheduler for student
	if err := es.db.Model(&e).Association("Schedulers").Append(sche); err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return e, nil
}

// func (es *EmployeeSvc) UpdateEmployeeScheduler(ctx context.Context, e *Employee, doorSerialId string, sche *Scheduler) (*Employee, error) {
// 	// Update scheduler for door
// 	var door = &Doorlock{}
// 	doorResult := es.db.Where("door_serial_id = ?", doorSerialId).First(door)
// 	if err := doorResult.Error; err != nil {
// 		err = utils.HandleQueryError(err)
// 		return nil, err
// 	}

// 	if err := es.db.Model(door).Association("Schedulers").Replace(sche); err != nil {
// 		err = utils.HandleQueryError(err)
// 		return nil, err
// 	}

// 	// Update scheduler for student
// 	if err := es.db.Model(&e).Association("Schedulers").Replace(sche); err != nil {
// 		err = utils.HandleQueryError(err)
// 		return nil, err
// 	}
// 	return e, nil
// }

// func (es *EmployeeSvc) DeleteEmployeeScheduler(ctx context.Context, e *Employee, doorSerialId string, sche *Scheduler) (*Employee, error) {
// 	// Delete scheduler for door
// 	var door = &Doorlock{}
// 	doorResult := es.db.Where("door_serial_id = ?", doorSerialId).First(door)
// 	if err := doorResult.Error; err != nil {
// 		err = utils.HandleQueryError(err)
// 		return nil, err
// 	}

// 	if err := es.db.Model(door).Association("Schedulers").Delete(sche); err != nil {
// 		err = utils.HandleQueryError(err)
// 		return nil, err
// 	}

// 	// Delete scheduler for student
// 	if err := es.db.Model(&e).Association("Schedulers").Delete(sche); err != nil {
// 		err = utils.HandleQueryError(err)
// 		return nil, err
// 	}
// 	return e, nil
// }
