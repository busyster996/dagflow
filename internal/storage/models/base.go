package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/busyster996/dagflow/internal/utility/sid"
)

type State int

const (
	StateStopped State = iota // 成功
	StateRunning              // 运行
	StateFailed               // 失败
	StateUnknown              // 未知
	StatePending              // 等待
	StatePaused               // 挂起
	StateSkipped              // 跳过
	StateAll     State = -1
)

var StateMap = map[State]string{
	StateStopped: "stopped",
	StateRunning: "running",
	StateFailed:  "failed",
	StateUnknown: "unknown",
	StatePending: "pending",
	StatePaused:  "paused",
	StateSkipped: "skipped",
}

type SBase struct {
	ID        uint64    `json:"id" gorm:"primarykey;comment:ID"`
	CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time `json:"updated_at" gorm:"comment:更新时间"`
}

func (b *SBase) BeforeCreate(tx *gorm.DB) (err error) {
	tableName := tx.Statement.Table
	b.ID, err = sid.NextID(tableName)
	if err != nil {
		return err
	}
	return
}

type SEnvs []*SEnv

type SEnv struct {
	Name  string `json:"name,omitempty" gorm:"size:256;index:,unique,composite:key;not null;comment:名称"`
	Value string `json:"value,omitempty" gorm:"size:256;comment:值"`
}
