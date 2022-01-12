package models

import (
	"context"
	"fmt"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type Student struct {
	GormModel
	MSSV  string `gorm:"type:varchar(256); unique; not null;" json:"mssv"  binding:"required"`
	Name  string `json:"name"`
	Phone string `gorm:"type:varchar(50)" json:"phone"`
	Email string `gorm:"type:varchar(256); unique; not null;" json:"email"`
	Major string `gorm:"not null;" json:"major"`
	UserPass
	Schedulers []Scheduler `gorm:"foreignKey:StudentID;references:MSSV;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"schedulers"`
}

type DeleteStudent struct {
	MSSV string `json:"mssv" binding:"required"`
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
	result := ss.db.Preload("Schedulers").Find(&sList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return sList, nil
}

func (ss *StudentSvc) FindStudentByMSSV(ctx context.Context, mssv string) (s *Student, err error) {
	var cnt int64
	result := ss.db.Preload("Schedulers").Where("mssv = ?", mssv).Find(&s).Count(&cnt)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	if cnt <= 0 {
		return nil, fmt.Errorf("find no records")
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
	result := ss.db.Model(&s).Where("id = ? AND mssv = ?", s.ID, s.MSSV).Updates(s)
	return utils.ReturnBoolStateFromResult(result)
}

func (ss *StudentSvc) DeleteStudent(ctx context.Context, mssv string) (bool, error) {
	result := ss.db.Unscoped().Where("mssv = ?", mssv).Delete(&Student{})
	return utils.ReturnBoolStateFromResult(result)
}

func (ss *StudentSvc) AppendStudentScheduler(ctx context.Context, s *Student, doorSerialId string, sche *Scheduler) (*Student, error) {

	// Add scheduler for door
	var door = &Doorlock{}
	doorResult := ss.db.Where("door_serial_id = ?", doorSerialId).First(door)
	if err := doorResult.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	sche.StudentID = &s.MSSV
	sche.DoorSerialID = &doorSerialId
	if err := ss.db.Create(&sche).Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	if err := ss.db.Model(door).Association("Schedulers").Append(sche); err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	// Add scheduler for student
	if err := ss.db.Model(&s).Association("Schedulers").Append(sche); err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return s, nil
}

// func (ss *StudentSvc) UpdateStudentScheduler(ctx context.Context, s *Student, doorSerialId string, sche *Scheduler) (*Student, error) {
// 	// Update scheduler for door
// 	var door = &Doorlock{}
// 	doorResult := ss.db.Where("door_serial_id = ?", doorSerialId).First(door)
// 	if err := doorResult.Error; err != nil {
// 		err = utils.HandleQueryError(err)
// 		return nil, err
// 	}

// 	if err := ss.db.Model(door).Association("Schedulers").Replace(sche); err != nil {
// 		err = utils.HandleQueryError(err)
// 		return nil, err
// 	}

// 	// Update scheduler for student
// 	if err := ss.db.Model(&s).Association("Schedulers").Replace(sche); err != nil {
// 		err = utils.HandleQueryError(err)
// 		return nil, err
// 	}
// 	return s, nil
// }

// func (ss *StudentSvc) DeleteStudentScheduler(ctx context.Context, s *Student, doorSerialId string, sche *Scheduler) (*Student, error) {
// 	// Delete scheduler for door
// 	var door = &Doorlock{}
// 	doorResult := ss.db.Where("door_serial_id = ?", doorSerialId).First(door)
// 	if err := doorResult.Error; err != nil {
// 		err = utils.HandleQueryError(err)
// 		return nil, err
// 	}

// 	if err := ss.db.Model(door).Association("Schedulers").Delete(sche); err != nil {
// 		err = utils.HandleQueryError(err)
// 		return nil, err
// 	}

// 	// Delete scheduler for student
// 	if err := ss.db.Model(&s).Association("Schedulers").Delete(sche); err != nil {
// 		err = utils.HandleQueryError(err)
// 		return nil, err
// 	}
// 	return s, nil
// }
