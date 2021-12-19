package models

import "time"

type GormModel struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `swaggerignore:"true"`
	UpdatedAt time.Time `swaggerignore:"true"`
}
