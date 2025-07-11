package models

import (
	"time"

	"gorm.io/datatypes"
)

type SRetryPolicy struct {
	Interval    time.Duration `json:"interval,omitempty" description:"间隔时间"`
	MaxInterval time.Duration `json:"maxInterval,omitempty" description:"最大间隔时间"`
	MaxAttempts int           `json:"maxAttempts,omitempty" description:"最大尝试次数"`
	Multiplier  float64       `json:"multiplier,omitempty" description:"乘数"`
}

type SStep struct {
	SBase
	TaskName    string                           `json:"task_name,omitempty" gorm:"size:256;uniqueIndex:idx_task_step_name;not null;comment:任务名称"`
	Name        string                           `json:"name,omitempty" gorm:"size:256;uniqueIndex:idx_task_step_name;not null;comment:名称"`
	Desc        string                           `json:"desc,omitempty" gorm:"comment:描述"`
	Type        string                           `json:"type,omitempty" gorm:"size:256;index;not null;comment:类型"`
	Content     string                           `json:"content,omitempty" gorm:"comment:内容"`
	Action      string                           `json:"action,omitempty" gorm:"comment:动作"`
	Rule        string                           `json:"rule,omitempty" gorm:"comment:规则"`
	RetryPolicy datatypes.JSONType[SRetryPolicy] `json:"retry_policy,omitempty" gorm:"comment:重试策略"`
	SeqNo       int64                            `json:"seq_no,omitempty" gorm:"index;not null;default:0;comment:序号"`
	Timeout     time.Duration                    `json:"timeout,omitempty" gorm:"not null;default:86400000000000;comment:超时时间"`
	Disable     *bool                            `json:"disable,omitempty" gorm:"not null;default:false;comment:禁用"`
	SStepUpdate
}

func (s *SStep) TableName() string {
	return "t_step"
}

type SStepUpdate struct {
	Message  string     `json:"message,omitempty" gorm:"comment:消息"`
	State    *State     `json:"state,omitempty" gorm:"index;not null;default:0;comment:状态"`
	OldState *State     `json:"old_state,omitempty" gorm:"index;not null;default:0;comment:旧状态"`
	Code     *int64     `json:"code,omitempty" gorm:"index;not null;default:0;comment:退出码"`
	STime    *time.Time `json:"s_time,omitempty" gorm:"comment:开始时间"`
	ETime    *time.Time `json:"e_time,omitempty" gorm:"comment:结束时间"`
}

func (s *SStepUpdate) STimeStr() string {
	if s.STime == nil {
		return "1970-01-01T00:00:00"
	}
	return s.STime.Format(time.RFC3339)
}

func (s *SStepUpdate) ETimeStr() string {
	if s.ETime == nil {
		return "1970-01-01T00:00:00"
	}
	return s.ETime.Format(time.RFC3339)
}

type SSteps []*SStep

func (s SSteps) Len() int {
	return len(s)
}

func (s SSteps) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

func (s SSteps) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
