package models

import (
	"context"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type Employee struct {
	GormModel
	MSNV       string `gorm:"type:varchar(50); unique; not null;" json:"msnv"`
	Name       string `gorm:"type:varchar(256)" json:"name"`
	Phone      string `gorm:"type:varchar(50)" json:"phone"`
	Email      string `gorm:"type:varchar(256); not null;" json:"email"`
	Department string `gorm:"type:varchar(50)" json:"department"`
	Role       string `gorm:"type:varchar(256); not null;" json:"role"`
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
	result := es.db.Find(&eList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return eList, nil
}

func (es *EmployeeSvc) FindEmployeeByID(ctx context.Context, id string) (e *Employee, err error) {
	result := es.db.First(&e, id)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return e, nil
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

func (es *EmployeeSvc) DeleteEmployee(ctx context.Context, employeeId uint) (bool, error) {
	err := es.db.Unscoped().Where("id = ?", employeeId).Delete(&Employee{}).Error
	if err != nil {
		return false, err
	}
	result := es.db.Unscoped().Where("id = ?", employeeId).Delete(&Password{})
	return utils.ReturnBoolStateFromResult(result)
}
