package models

import (
	"context"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type Student struct {
	GormModel
	MSSV       string      `gorm:"type:varchar(50); unique; not null;" json:"mssv"`
	Name       string      `json:"name"`
	Phone      string      `gorm:"type:varchar(50)" json:"phone"`
	Email      string      `gorm:"type:varchar(256); unique; not null;" json:"email"`
	Major      string      `gorm:"not null;" json:"major"`
	Schedulers []Scheduler `gorm:"many2many:student_schedulers;"`
}

type StudentSchedulerUpsert struct {
	Scheduler  `json:"scheduler" binding:"required"`
	GatewayID  string `json:"gatewayId" binding:"required"`
	DoorlockID string `json:"doorlockId" binding:"required"`
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

func (ss *StudentSvc) FindStudentByID(ctx context.Context, id string) (s *Student, err error) {
	result := ss.db.Preload("Schedulers").First(&s, id)
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
	result := ss.db.Unscoped().Where("id = ?", studentId).Delete(&Student{})
	isSuccess, err := utils.ReturnBoolStateFromResult(result)
	if !isSuccess {
		return isSuccess, err
	}
	err = ss.db.Unscoped().Where("id = ?", studentId).Delete(&Password{}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (ss *StudentSvc) AppendStudentScheduler(ctx context.Context, sId string, doorSerialId string, sche *Scheduler) (*Student, error) {

	// Add scheduler for door
	var door = &Doorlock{}
	doorResult := ss.db.Where("door_serial_id = ?", doorSerialId).First(door)
	if err := doorResult.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	if err := ss.db.Model(door).Association("Schedulers").Append(sche); err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	// Add scheduler for student
	s, err := ss.FindStudentByID(ctx, sId)
	if err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	if err := ss.db.Model(&s).Association("Schedulers").Append(sche); err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return s, nil
}

func (ss *StudentSvc) UpdateStudentScheduler(ctx context.Context, sId string, doorSerialId string, sche *Scheduler) (*Student, error) {
	// Update scheduler for door
	var door = &Doorlock{}
	doorResult := ss.db.Where("door_serial_id = ?", doorSerialId).First(door)
	if err := doorResult.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	if err := ss.db.Model(door).Association("Schedulers").Replace(sche); err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	// Update scheduler for student
	s, err := ss.FindStudentByID(ctx, sId)
	if err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	if err := ss.db.Model(&s).Association("Schedulers").Replace(sche); err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return s, nil
}

func (ss *StudentSvc) DeleteStudentScheduler(ctx context.Context, sId string, doorSerialId string, sche *Scheduler) (*Student, error) {
	// Delete scheduler for door
	var door = &Doorlock{}
	doorResult := ss.db.Where("door_serial_id = ?", doorSerialId).First(door)
	if err := doorResult.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	if err := ss.db.Model(door).Association("Schedulers").Delete(sche); err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	// Delete scheduler for student
	s, err := ss.FindStudentByID(ctx, sId)
	if err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}

	if err := ss.db.Model(&s).Association("Schedulers").Delete(sche); err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return s, nil
}
