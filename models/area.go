package models

import (
	"context"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type Area struct {
	gorm.Model
	Gateway   Gateway    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Name      string     `gorm:"unique;not null" json:"name"`
	Manager   string     `gorm:"unique;not null" json:"manager"`
	Doorlocks []Doorlock `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
type AreaSvc struct {
	db *gorm.DB
}

func NewAreaSvc(db *gorm.DB) *AreaSvc {
	return &AreaSvc{
		db: db,
	}
}

func (as *AreaSvc) FindAllArea(ctx context.Context) (aList []Area, err error) {
	result := as.db.Joins("Gateway").Find(&aList)
	if err := result.Error; err != nil {
		err = utils.QueryErrorHandler(err)
		return nil, err
	}
	return aList, nil
}

func (as *AreaSvc) FindAreaByID(ctx context.Context, id uint) (a *Area, err error) {
	result := as.db.First(&a, id)
	if err := result.Error; err != nil {
		err = utils.QueryErrorHandler(err)
		return nil, err
	}
	return a, nil
}

func (as *AreaSvc) CreateArea(a *Area, ctx context.Context) (*Area, error) {
	if err := as.db.Create(&a).Error; err != nil {
		err = utils.QueryErrorHandler(err)
		return nil, err
	}
	return a, nil
}

func (as *AreaSvc) UpdateArea(ctx context.Context, a *Area) (*Area, error) {
	result := as.db.Model(&a).Where("id = ?", a.ID).Updates(a)
	handled, err := utils.UpdateResultHandler(result, a)
	if err != nil {
		return nil, err
	}
	a = handled.(*Area)
	return a, nil
}

func (as *AreaSvc) DeleteArea(ctx context.Context, a *Area) (bool, error) {
	result := as.db.Unscoped().Where("id = ?", a.ID).Delete(a)
	return utils.DeleteResultHandler(result)
}
