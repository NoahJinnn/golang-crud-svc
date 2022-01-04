package models

import (
	"context"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type Student struct {
	GormModel
	MSSV  string `gorm:"type:varchar(50); unique; not null;" json:"mssv"`
	Name  string `gorm:"type:varchar(256)" json:"name"`
	Phone string `gorm:"type:varchar(50)" json:"phone"`
	Email string `gorm:"type:varchar(256); unique; not null;" json:"email"`
	Major string `gorm:"type:varchar(256); not null;" json:"major"`
}

type StudentSvc struct {
	db *gorm.DB
}

func NewStudentSvc(db *gorm.DB) *StudentSvc {
	return &StudentSvc{
		db: db,
	}
}

func (ss *StudentSvc) FindAllStudent(ctx context.Context) (sList []Student, err error) {
	result := ss.db.Find(&sList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return sList, nil
}

func (ss *StudentSvc) FindStudentByID(ctx context.Context, id string) (s *Student, err error) {
	result := ss.db.First(&s, id)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return s, nil
}

func (ss *StudentSvc) CreateStudent(ctx context.Context, s *Student) (*Student, error) {
	if err := ss.db.Create(&s).Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return s, nil
}

func (ss *StudentSvc) UpdateStudent(ctx context.Context, s *Student) (bool, error) {
	result := ss.db.Model(&s).Where("id = ?", s.ID).Updates(s)
	return utils.ReturnBoolStateFromResult(result)
}

func (ss *StudentSvc) DeleteStudent(ctx context.Context, studentId uint) (bool, error) {
	err := ss.db.Unscoped().Where("id = ?", studentId).Delete(&Student{}).Error
	if err != nil {
		return false, err
	}
	result := ss.db.Unscoped().Where("id = ?", studentId).Delete(&Password{})
	return utils.ReturnBoolStateFromResult(result)
}
