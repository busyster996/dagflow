package model

type TaskParam struct {
	Base
	TaskID  uint64 `gorm:"index:,unique,composite:key;not null;comment:任务ID"`
	ParamID uint64 `gorm:"index:,unique,composite:key;not null;comment:参数ID"`
}

func (t *TaskParam) TableName() string {
	return "task_params"
}
