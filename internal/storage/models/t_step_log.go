package models

type SStepLog struct {
	SBase
	TaskName  string `json:"task_name,omitempty" gorm:"size:256;index;not null;comment:任务名称"`
	StepName  string `json:"step_name,omitempty" gorm:"size:256;index;not null;comment:步骤名称"`
	Timestamp int64  `json:"timestamp,omitempty" gorm:"not null;comment:时间戳"`
	Line      *int64 `json:"line,omitempty" gorm:"not null;comment:行号"`
	Content   string `json:"content,omitempty" gorm:"comment:内容"`
}

func (l *SStepLog) TableName() string {
	return "t_step_log"
}

type SStepLogs []*SStepLog

func (l SStepLogs) Len() int {
	return len(l)
}

func (l SStepLogs) Less(i, j int) bool {
	if l[i].Line == nil || l[j].Line == nil {
		return l[i].Timestamp < l[j].Timestamp
	}
	return *l[i].Line < *l[j].Line
}

func (l SStepLogs) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
