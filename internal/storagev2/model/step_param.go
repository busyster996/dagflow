package model

type StepParam struct {
	Base
	StepID  uint64 `gorm:"index:,unique,composite:key;not null;comment:步骤ID"`
	ParamID uint64 `gorm:"index:,unique,composite:key;not null;comment:参数ID"`
}

func (s *StepParam) TableName() string {
	return "step_params"
}
