package models

import (
	"gorm.io/gorm"
)

type Area struct {
	gorm.Model
	Gateway Gateway
	Name    string `gorm:unique;not null json:"name"`
}
