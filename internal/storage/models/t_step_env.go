package models

type SStepEnv struct {
	SBase
	TaskName string `json:"task_name,omitempty" gorm:"size:256;index:,unique,composite:key;not null;comment:任务名称"`
	StepName string `json:"step_name,omitempty" gorm:"size:256;index:,unique,composite:key;not null;comment:步骤名称"`
	SEnv
}

func (s *SStepEnv) TableName() string {
	return "t_step_env"
}
