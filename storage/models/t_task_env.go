package models

type STaskEnv struct {
	SBase
	TaskName string `json:"task_name,omitempty" gorm:"size:256;index:,unique,composite:key;not null;comment:任务名称"`
	SEnv
}

func (t *STaskEnv) TableName() string {
	return "t_task_env"
}
