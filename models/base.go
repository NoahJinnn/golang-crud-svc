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

type UserPass struct {
	RfidPass   string `gorm:"type:varchar(256)" json:"rfidPass"`
	KeypadPass string `gorm:"type:varchar(256)" json:"keypadPass"`
}

type UserSchedulerUpsert struct {
	Scheduler  `json:"scheduler" binding:"required"`
	GatewayID  string `json:"gatewayId" binding:"required"`
	DoorlockID string `json:"doorlockId" binding:"required"`
}
