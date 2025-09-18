package model

import (
	"time"
)

type Execution struct {
	Base
	TaskID        uint64     `gorm:"index:,unique,composite:key;not null;comment:任务ID"`
	StepID        uint64     `gorm:"index:,unique,composite:key;not null;comment:步骤ID"`
	Message       string     `gorm:"comment:消息"`
	ExitCode      *int64     `gorm:"index;not null;default:0;comment:退出码"`
	State         *ExitState `gorm:"index;not null;default:4;comment:状态"`
	PreviousState *ExitState `gorm:"index;not null;default:4;comment:旧状态"`
	StartTime     *time.Time `gorm:"comment:开始时间"`
	EndTime       *time.Time `gorm:"comment:结束时间"`

	// 关联关系
	Task    *Task     `gorm:"foreignKey:TaskID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Step    *Step     `gorm:"foreignKey:StepID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Outputs []*Output `gorm:"foreignKey:ExecutionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (e *Execution) TableName() string {
	return "executions"
}
