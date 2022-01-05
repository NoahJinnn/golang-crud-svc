package models

import (
	"context"

	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"gorm.io/gorm"
)

type Scheduler struct {
	GormModel
	Base           string `gorm:"not null;" json:"base"`
	RoomRow        string `gorm:"not null;" json:"roomRow"`
	RoomID         string `gorm:" not null;" json:"roomId"`
	RoomName       string `gorm:" not null;" json:"roomName"`
	StartDate      string `gorm:"type:varchar(50) not null;" json:"startDate"`
	EndDate        string `gorm:"type:varchar(50) not null;" json:"endDate"`
	ClassID        string `gorm:" not null;" json:"classId"`
	ClassName      string `gorm:" not null;" json:"className"`
	LecturerID     string `gorm:" not null;" json:"lecturerId"`
	LecturerName   string `gorm:" not null;" json:"lecturerName"`
	Capacity       uint   `json:"capacity"`
	WeekDay        uint   `json:"weekDay"`
	StartClassTime uint   `json:"startClassTime"`
	EndClassTime   uint   `json:"endClassTime"`
	Amount         uint   `json:"amount"`
	Status         string `json:"status"`
}

type SchedulerSvc struct {
	db *gorm.DB
}

func NewSchedulerSvc(db *gorm.DB) *SchedulerSvc {
	return &SchedulerSvc{
		db: db,
	}
}

func (ss *SchedulerSvc) FindAllScheduler(ctx context.Context) (sList []Scheduler, err error) {
	result := ss.db.Find(&sList)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return sList, nil
}

func (ss *SchedulerSvc) FindSchedulerByID(ctx context.Context, id string) (s *Scheduler, err error) {
	result := ss.db.First(&s, id)
	if err := result.Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return s, nil
}

func (ss *SchedulerSvc) CreateScheduler(ctx context.Context, s *Scheduler) (*Scheduler, error) {
	if err := ss.db.Create(&s).Error; err != nil {
		err = utils.HandleQueryError(err)
		return nil, err
	}
	return s, nil
}

func (ss *SchedulerSvc) UpdateScheduler(ctx context.Context, s *Scheduler) (bool, error) {
	result := ss.db.Model(&s).Where("id = ?", s.ID).Updates(s)
	return utils.ReturnBoolStateFromResult(result)
}

func (ss *SchedulerSvc) DeleteScheduler(ctx context.Context, studentId uint) (bool, error) {
	err := ss.db.Unscoped().Where("id = ?", studentId).Delete(&Scheduler{}).Error
	if err != nil {
		return false, err
	}
	result := ss.db.Unscoped().Where("id = ?", studentId).Delete(&Password{})
	return utils.ReturnBoolStateFromResult(result)
}
