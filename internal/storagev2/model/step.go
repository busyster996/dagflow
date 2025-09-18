package model

import (
	"time"

	"gorm.io/datatypes"
)

type Step struct {
	Base
	Name        string                          `gorm:"index;not null;comment:名称"`
	Desc        string                          `gorm:"comment:描述"`
	Kind        string                          `gorm:"index;not null;comment:类型"`
	Content     string                          `gorm:"comment:内容"`
	Action      string                          `gorm:"comment:动作"`
	Rule        string                          `gorm:"comment:规则"`
	Timeout     time.Duration                   `gorm:"not null;default:86400000000000;comment:超时时间"`
	Disabled    *bool                           `gorm:"not null;default:false;comment:禁用"`
	RetryPolicy datatypes.JSONType[RetryPolicy] `gorm:"comment:重试策略"`
	Metadata    datatypes.JSONMap               `gorm:"comment:元数据"`

	// 步骤依赖关系 - 自关联多对多
	Dependencies []*Step `gorm:"many2many:step_depends;joinForeignKey:StepID;joinReferences:DependID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Dependents   []*Step `gorm:"many2many:step_depends;joinForeignKey:DependID;joinReferences:StepID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// 关联关系
	Executions []*Execution `gorm:"foreignKey:StepID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Params     []*Param     `gorm:"many2many:step_params;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (s *Step) TableName() string {
	return "steps"
}
