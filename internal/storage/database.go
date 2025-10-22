package storage

import (
	"fmt"
	"runtime/debug"

	"gorm.io/gorm"

	"github.com/busyster996/dagflow/internal/common"
	"github.com/busyster996/dagflow/internal/storage/models"
	"github.com/busyster996/dagflow/pkg/logx"
)

type sDatabase struct {
	*gorm.DB
}

func (d *sDatabase) initSqlite() {
	d.Exec("PRAGMA mode=rwc;")
	// 开启外键约束
	d.Exec("PRAGMA foreign_keys=ON;")
	// 写同步
	d.Exec("PRAGMA synchronous=NORMAL;")
	// 启用 WAL 模式
	d.Exec("PRAGMA journal_mode=WAL;")
	// 控制WAL文件大小 100MB
	d.Exec("PRAGMA journal_size_limit=104857600;")
	// 设置等待超时，减少锁等待时间 5秒
	d.Exec("PRAGMA busy_timeout=5000;")
	// 设置共享缓存
	d.Exec("PRAGMA cache=shared;")
	// 设置缓存大小 约32MB缓存
	d.Exec("PRAGMA cache_size=-8000;")
	// 设置内存映射大小 128MB
	d.Exec("PRAGMA mmap_size=134217728;")
	// 将临时表放入内存
	d.Exec("PRAGMA temp_store=MEMORY;")
	// 设置锁模式为NORMAL，支持高并发访问
	d.Exec("PRAGMA locking_mode=NORMAL;")
	// 开启缓存溢出管理，适用于高并发写入
	d.Exec("PRAGMA cache_spill=ON;")
}

func (d *sDatabase) Name() string {
	return d.DB.Name()
}

func (d *sDatabase) FixDatabase(node string) (err error) {
	// 开始事务
	tx := d.Begin()
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			logx.Errorln(r, string(stack))
			tx.Rollback()
		}
	}()

	// 更新所有符合条件的步骤状态为失败
	if err = tx.Model(&models.SStep{}).
		Where("task_name IN (?)",
			d.Model(&models.STask{}).Select("name").
				Where("(node IS NULL OR node = ?) AND (state <> ? AND state <> ? AND state <> ?)", node, models.StateStopped, models.StateSkipped, models.StateFailed),
		).
		Where("state = ? OR state = ?", models.StateRunning, models.StatePaused).
		Updates(map[string]interface{}{
			"state":   models.StateFailed,
			"code":    common.ExecCodeSystemErr,
			"message": "execution failed due to system error",
		}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新所有符合条件的任务状态为失败
	if err = tx.Model(&models.STask{}).
		Where("(node IS NULL OR node = ?) AND (state <> ? AND state <> ? AND state <> ?)", node, models.StateStopped, models.StateSkipped, models.StateFailed).
		Updates(map[string]interface{}{
			"node":    node,
			"state":   models.StateFailed,
			"message": "execution failed due to system error",
		}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		return err
	}
	return
}

func (d *sDatabase) NodeTasks(node string) (res []ITask) {
	var tasks []string
	d.Model(&models.STask{}).Where("node = ?", node).Pluck("name", &tasks)
	for _, name := range tasks {
		res = append(res, &sTask{
			DB:    d.DB,
			tName: name,
		})
	}
	return
}

func (d *sDatabase) Task(name string) ITask {
	return &sTask{
		DB:    d.DB,
		tName: name,
	}
}

func (d *sDatabase) TaskCreate(task *models.STask) (err error) {
	if task.Node == "random" {
		task.Node = ""
	}
	err = d.Create(task).Error
	if err != nil {
		return fmt.Errorf("save task %d error: %s", task.ID, err)
	}
	return
}

func (d *sDatabase) TaskCount(state models.State) (res int64) {
	if state != models.StateAll {
		d.Model(&models.STask{}).Distinct("DISTINCT name").Where("state = ?", state).Count(&res)
		return
	}
	d.Model(&models.STask{}).Distinct("DISTINCT name").Count(&res)
	return
}

func (d *sDatabase) TaskList(page, pageSize int64, str string) (res models.STasks, total int64) {
	err := d.Model(&models.STask{}).Count(&total).Error
	if err != nil {
		return
	}
	query := d.Model(&models.STask{}).
		Select("id, name, kind, state, message, s_time, e_time").
		Order("id DESC")
	if str != "" {
		query.Where("name LIKE ?", str+"%")
	}
	query.Scopes(func(db *gorm.DB) *gorm.DB {
		return models.Paginate(db, page, pageSize)
	}).Find(&res)
	return
}

func (d *sDatabase) Pipeline(name string) IPipeline {
	return &sPipeline{
		DB:   d.DB,
		name: name,
	}
}

func (d *sDatabase) PipelineCreate(pipeline *models.SPipeline) (err error) {
	err = d.Create(pipeline).Error
	if err != nil {
		return fmt.Errorf("save pipeline %d error: %s", pipeline.ID, err)
	}
	return
}

func (d *sDatabase) PipelineList(page, pageSize int64, str string) (res models.SPipelines, total int64) {
	err := d.Model(&models.SPipeline{}).Count(&total).Error
	if err != nil {
		return
	}
	query := d.Model(&models.SPipeline{}).
		Select("id, name, disable, tpl_type").
		Order("id DESC")
	if str != "" {
		query.Where("name LIKE ?", str+"%")
	}
	query.Scopes(func(db *gorm.DB) *gorm.DB {
		return models.Paginate(db, page, pageSize)
	}).Find(&res)
	return
}
