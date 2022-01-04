package models

import (
	"time"
)

type GormModel struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `swaggerignore:"true"`
	UpdatedAt time.Time `swaggerignore:"true"`
}

type DeleteID struct {
	ID uint `json:"id"`
}
