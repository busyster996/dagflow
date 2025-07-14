package storage

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/busyster996/dagflow/internal/storage/models"
	"github.com/busyster996/dagflow/pkg/logx"
)

type sStepLog struct {
	*gorm.DB
	tName string
	sName string

	lock sync.Mutex
}

func (l *sStepLog) List(latestLine *int64) (res models.SStepLogs) {
	query := l.Model(&models.SStepLog{}).
		Where(map[string]interface{}{
			"task_name": l.tName,
			"step_name": l.sName,
		}).Order("line ASC")
	if latestLine != nil {
		// 如果 latestLine 不为空，只查询行号大于 latestLine 的日志
		query = query.Where("line > ?", latestLine).Limit(500)
	}
	query.Find(&res)
	return
}

func (l *sStepLog) Insert(log *models.SStepLog) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	log.TaskName = l.tName
	log.StepName = l.sName
	return l.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := l.Model(&models.SStepLog{}).
			Where(map[string]interface{}{
				"task_name": l.tName,
				"step_name": l.sName,
			}).
			Count(&count).Error; err != nil {
			return err
		}
		log.Line = models.Pointer(count)
		return tx.Create(log).Error
	})
}

func (l *sStepLog) Write(contents ...string) {
	if err := l.Insert(&models.SStepLog{
		Timestamp: time.Now().UnixNano(),
		Content:   strings.Join(contents, " "),
	}); err != nil {
		logx.Warnln(err)
	}
}

func (l *sStepLog) Writef(format string, args ...interface{}) {
	l.Write(fmt.Sprintf(format, args...))
}

func (l *sStepLog) RemoveAll() (err error) {
	return l.Where(map[string]interface{}{
		"task_name": l.tName,
		"step_name": l.sName,
	}).Delete(&models.SStepLog{}).Error
}
