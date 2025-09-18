package model

import (
	"time"

	"gorm.io/datatypes"
)

type Task struct {
	Base
	Kind        string                          `gorm:"index;comment:类型"`
	Name        string                          `gorm:"not null;comment:名称"`
	Desc        string                          `gorm:"comment:描述"`
	Node        string                          `gorm:"index;default:null;comment:节点"`
	Timeout     time.Duration                   `gorm:"not null;default:86400000000000;comment:超时时间"`
	Disabled    *bool                           `gorm:"not null;default:false;comment:禁用"`
	IsTpl       *bool                           `gorm:"not null;default:false;comment:是否模板"`
	RetryPolicy datatypes.JSONType[RetryPolicy] `gorm:"comment:重试策略"`
	Metadata    datatypes.JSONMap               `gorm:"comment:元数据"`
	Message     string                          `gorm:"comment:消息"`
	State       *ExitState                      `gorm:"index;not null;default:4;comment:状态"`
	OldState    *ExitState                      `gorm:"index;not null;default:4;comment:旧状态"`
	StartTime   *time.Time                      `gorm:"comment:开始时间"`
	EndTime     *time.Time                      `gorm:"comment:结束时间"`

	// 关联关系
	Executions []*Execution `gorm:"foreignKey:TaskID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Params     []*Param     `gorm:"many2many:task_params;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (t *Task) TableName() string {
	return "tasks"
}
