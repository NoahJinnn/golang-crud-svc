package models

import (
	"context"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type Area struct {
	gorm.Model
	Gateway Gateway `gorm:"foreignKey:AreaID; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Name    string  `gorm:"unique;not null" json:"name"`
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

func (as *AreaSvc) DeleteArea(ctx context.Context, a *Area) (bool, error) {
	result := as.db.Unscoped().Where("id = ?", a.ID).Delete(a)
	return utils.DeleteResultHandler(result)
}
