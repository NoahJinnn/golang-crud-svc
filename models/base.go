package models

import (
	"time"
)

type GormModel struct {
	ID        string    `gorm:"primarykey; type:varchar(256)"`
	CreatedAt time.Time `swaggerignore:"true"`
	UpdatedAt time.Time `swaggerignore:"true"`
}
