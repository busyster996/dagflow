package model

type Param struct {
	Base
	Name  string  `gorm:"index:,unique,composite:key;not null;comment:名称"`
	Desc  string  `gorm:"comment:描述"`
	Value *string `gorm:"comment:值"`

	// 关联关系
	Tasks []*Task `gorm:"many2many:task_params;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Steps []*Step `gorm:"many2many:step_params;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (p *Param) TableName() string {
	return "params"
}
