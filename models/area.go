package models

import (
	"context"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Area struct {
	GormModel
	Gateway   Gateway    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"gateway"`
	Name      string     `gorm:"unique;not null" json:"name"`
	Manager   string     `gorm:"unique;not null" json:"manager"`
	Doorlocks []Doorlock `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"doorlocks"`
}
type AreaSvc struct {
	db *gorm.DB
}

func NewAreaSvc(db *gorm.DB) *AreaSvc {
	return &AreaSvc{
		db: db,
	}
}

func (a *Area) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.New().String()
	return
}

func (as *AreaSvc) FindAllArea(ctx context.Context) (aList []Area, err error) {
	result := as.db.Preload("Doorlocks").Joins("Gateway").Find(&aList)
	if err := result.Error; err != nil {
		err = utils.QueryErrorHandler(err)
		return nil, err
	}
	return aList, nil
}

func (as *AreaSvc) FindAreaByID(ctx context.Context, id string) (a *Area, err error) {
	result := as.db.Preload("Doorlocks").First(&a, id)
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

func (as *AreaSvc) UpdateArea(ctx context.Context, a *Area) (bool, error) {
	result := as.db.Model(&a).Where("id = ?", a.ID).Updates(a)
	return utils.ReturnBoolStateFromResult(result)
}

func (as *AreaSvc) DeleteArea(ctx context.Context, a *Area) (bool, error) {
	result := as.db.Unscoped().Where("id = ?", a.ID).Delete(a)
	return utils.ReturnBoolStateFromResult(result)
}
