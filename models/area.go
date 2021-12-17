package models

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Area struct {
	gorm.Model
	Gateway Gateway
	Name    string `gorm:"unique;not null" json:"name"`
}
type AreaSvc struct {
	db *gorm.DB
}

func NewAreaSvc(db *gorm.DB) *AreaSvc {
	return &AreaSvc{
		db: db,
	}
}

func (as *AreaSvc) CreateArea(a *Area, ctx context.Context) (*Area, error) {
	if err := as.db.Create(&a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return a, nil
}
