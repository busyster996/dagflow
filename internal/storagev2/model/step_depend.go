package model

type StepDepend struct {
	Base
	StepID   uint64 `gorm:"index:,unique,composite:key;not null;comment:步骤ID"`
	DependID uint64 `gorm:"index:,unique,composite:key;not null;comment:依赖步骤ID"`
}

func (s *StepDepend) TableName() string {
	return "step_depends"
}
