package model

import (
	"time"

	"gorm.io/gorm"

	"github.com/busyster996/dagflow/internal/utils/sid"
)

type ExitState int

const (
	ExitStateStopped ExitState = iota // 成功
	ExitStateRunning                  // 运行
	ExitStateFailed                   // 失败
	ExitStateUnknown                  // 未知
	ExitStatePending                  // 等待
	ExitStatePaused                   // 挂起
	ExitStateSkipped                  // 跳过
)

var (
	ExitStateMap = map[ExitState]string{
		ExitStateStopped: "stopped",
		ExitStateRunning: "running",
		ExitStateFailed:  "failed",
		ExitStateUnknown: "unknown",
		ExitStatePending: "pending",
		ExitStatePaused:  "paused",
		ExitStateSkipped: "skipped",
	}
)

type Base struct {
	ID        uint64    `gorm:"primarykey;comment:ID"`
	CreatedAt time.Time `gorm:"comment:创建时间"`
	UpdatedAt time.Time `gorm:"comment:更新时间"`
	CreatedBy *string   `gorm:"comment:创建者"`
	UpdatedBy *string   `gorm:"comment:更新者"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	tableName := tx.Statement.Table
	b.ID, err = sid.NextID(tableName)
	return
}

type RetryPolicy struct {
	Interval    time.Duration `json:"interval,omitempty" yaml:"interval,omitempty" description:"间隔时间"`
	MaxInterval time.Duration `json:"maxInterval,omitempty" yaml:"maxInterval,omitempty" description:"最大间隔时间"`
	MaxAttempts int           `json:"maxAttempts,omitempty" yaml:"maxAttempts,omitempty" description:"最大尝试次数"`
	Multiplier  float64       `json:"multiplier,omitempty" yaml:"multiplier,omitempty" description:"乘数"`
}

func Pointer[T any](v T) *T {
	return &v
}

func Paginate(db *gorm.DB, page, pageSize int64) *gorm.DB {
	if page == -1 {
		return db
	}
	if page == 0 {
		page = 1
	}
	switch {
	case pageSize > 500:
		pageSize = 500
	case pageSize <= 0:
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	return db.Offset(int(offset)).Limit(int(pageSize))
}
