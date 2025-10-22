package models

import (
	"time"

	"gorm.io/datatypes"
)

type STask struct {
	SBase
	Kind     string            `json:"kind,omitempty" gorm:"size:256;index;comment:类型"`
	Name     string            `json:"name,omitempty" gorm:"size:256;uniqueIndex;not null;comment:名称"`
	Desc     string            `json:"desc,omitempty" gorm:"comment:描述"`
	Node     string            `json:"node,omitempty" gorm:"size:256;index;default:null;comment:节点"`
	Timeout  time.Duration     `json:"timeout,omitempty" gorm:"not null;default:86400000000000;comment:超时时间"`
	Disable  *bool             `json:"disable,omitempty" gorm:"not null;default:false;comment:禁用"`
	Metadata datatypes.JSONMap `json:"metadata,omitempty" gorm:"元数据"`
	STaskUpdate
}

func (t *STask) TableName() string {
	return "t_task"
}

type STaskUpdate struct {
	Message  string     `json:"message,omitempty" gorm:"comment:消息"`
	State    *State     `json:"state,omitempty" gorm:"index;not null;default:0;comment:状态"`
	OldState *State     `json:"old_state,omitempty" gorm:"index;not null;default:0;comment:旧状态"`
	STime    *time.Time `json:"s_time,omitempty" gorm:"comment:开始时间"`
	ETime    *time.Time `json:"e_time,omitempty" gorm:"comment:结束时间"`
}

func (t *STaskUpdate) STimeStr() string {
	if t.STime == nil {
		return "1970-01-01T00:00:00"
	}
	return t.STime.Format(time.RFC3339)
}

func (t *STaskUpdate) ETimeStr() string {
	if t.ETime == nil {
		return "1970-01-01T00:00:00"
	}
	return t.ETime.Format(time.RFC3339)
}

type STasks []*STask

func (t STasks) Len() int {
	return len(t)
}

func (t STasks) Less(i, j int) bool {
	if t[i].STime == nil || t[j].STime == nil {
		return t[i].Name < t[j].Name
	}
	iTime := t[i].STime.Nanosecond()
	jTime := t[j].STime.Nanosecond()
	return iTime < jTime
}

func (t STasks) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
