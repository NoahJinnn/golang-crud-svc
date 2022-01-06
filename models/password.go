package models

// import (
// 	"context"

// 	"github.com/ecoprohcm/DMS_BackendServer/utils"
// 	"gorm.io/gorm"
// )

// type Password struct {
// 	GormModel
// 	UserID       uint   `gorm:"not null;" json:"userId"`
// 	PasswordType string `gorm:"type:varchar(10)" json:"passwordType"`
// 	PasswordHash string `gorm:"type:varchar(256)" json:"passwordHash"`
// }

// type PasswordCreate struct {
// 	Password
// 	GatewayID string `json:"gatewayId"`
// }
// type PasswordSvc struct {
// 	db *gorm.DB
// }

// func NewPasswordSvc(db *gorm.DB) *PasswordSvc {
// 	return &PasswordSvc{
// 		db: db,
// 	}
// }

// func (ps *PasswordSvc) FindAllPasswordByUserID(ctx context.Context, userId string) (pList []Password, err error) {
// 	result := ps.db.Where("user_id = ?", userId).Find(&pList)
// 	if err := result.Error; err != nil {
// 		err = utils.HandleQueryError(err)
// 		return nil, err
// 	}
// 	return pList, nil
// }

// func (ps *PasswordSvc) CreatePassword(ctx context.Context, p *Password) (*Password, error) {
// 	if err := ps.db.Create(&p).Error; err != nil {
// 		err = utils.HandleQueryError(err)
// 		return nil, err
// 	}
// 	return p, nil
// }

// func (ps *PasswordSvc) UpdatePassword(ctx context.Context, p *Password) (bool, error) {
// 	result := ps.db.Model(&p).Where("id = ?", p.ID).Updates(p)
// 	return utils.ReturnBoolStateFromResult(result)
// }

// func (ps *PasswordSvc) DeletePassword(ctx context.Context, pwId uint) (bool, error) {
// 	result := ps.db.Unscoped().Where("id = ?", pwId).Delete(&Password{})
// 	return utils.ReturnBoolStateFromResult(result)
// }
