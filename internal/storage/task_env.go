package storage

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/busyster996/dagflow/internal/storage/models"
)

type sTaskEnv struct {
	*gorm.DB
	tName string
}

func (e *sTaskEnv) List() (res models.SEnvs) {
	e.Model(&models.STaskEnv{}).
		Where(map[string]interface{}{
			"task_name": e.tName,
		}).
		Order("id ASC").
		Find(&res)
	return
}

func (e *sTaskEnv) Insert(envs ...*models.SEnv) (err error) {
	if len(envs) == 0 {
		return
	}
	var _envs []models.STaskEnv
	for _, env := range envs {
		_envs = append(_envs, models.STaskEnv{
			TaskName: e.tName,
			SEnv:     *env,
		})
	}
	err = e.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "task_name"},
			{Name: "name"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(_envs).Error
	if err != nil {
		return err
	}
	return
}

func (e *sTaskEnv) Update(env *models.SEnv) (err error) {
	return e.Model(&models.STaskEnv{}).
		Omit("name").
		Where(map[string]interface{}{
			"task_name": e.tName,
			"name":      env.Name,
		}).
		Updates(env).Error
}

func (e *sTaskEnv) Get(name string) (res string, err error) {
	if name == "" {
		return "", errors.New("name is empty")
	}
	err = e.Model(&models.STaskEnv{}).
		Select("value").
		Where(map[string]interface{}{
			"task_name": e.tName,
			"name":      name,
		}).
		Scan(&res).
		Error
	return
}

func (e *sTaskEnv) Remove(name string) (err error) {
	if name == "" {
		return errors.New("name is empty")
	}
	return e.Where(map[string]interface{}{
		"task_name": e.tName,
		"name":      name,
	}).Delete(&models.STaskEnv{}).Error
}

func (e *sTaskEnv) RemoveAll() (err error) {
	return e.Where(map[string]interface{}{
		"task_name": e.tName,
	}).Delete(&models.STaskEnv{}).Error
}
