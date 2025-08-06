package storage

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/busyster996/dagflow/internal/storage/models"
)

type sStepEnv struct {
	*gorm.DB
	tName string
	sName string
}

func (e *sStepEnv) List() (res models.SEnvs) {
	e.Model(&models.SStepEnv{}).
		Where(map[string]interface{}{
			"task_name": e.tName,
			"step_name": e.sName,
		}).
		Order("id ASC").
		Find(&res)
	return
}

func (e *sStepEnv) Insert(envs ...*models.SEnv) (err error) {
	if len(envs) == 0 {
		return
	}
	var _envs []models.SStepEnv
	for _, env := range envs {
		_envs = append(_envs, models.SStepEnv{
			TaskName: e.tName,
			StepName: e.sName,
			SEnv:     *env,
		})
	}
	err = e.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "task_name"},
			{Name: "step_name"},
			{Name: "name"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(_envs).Error
	if err != nil {
		return err
	}
	return
}

func (e *sStepEnv) Update(env *models.SEnv) (err error) {
	return e.Model(&models.SStepEnv{}).
		Omit("name").
		Where(map[string]interface{}{
			"task_name": e.tName,
			"step_name": e.sName,
			"name":      env.Name,
		}).
		Updates(env).Error
}

func (e *sStepEnv) Get(name string) (res string, err error) {
	if name == "" {
		return "", errors.New("name is empty")
	}
	err = e.Model(&models.SStepEnv{}).
		Select("value").
		Where(map[string]interface{}{
			"task_name": e.tName,
			"step_name": e.sName,
			"name":      name,
		}).
		Scan(&res).
		Error
	return
}

func (e *sStepEnv) Remove(name string) (err error) {
	if name == "" {
		return errors.New("name is empty")
	}
	return e.Where(map[string]interface{}{
		"task_name": e.tName,
		"step_name": e.sName,
		"name":      name,
	}).Delete(&models.SStepEnv{}).Error
}

func (e *sStepEnv) RemoveAll() (err error) {
	return e.Where(map[string]interface{}{
		"task_name": e.tName,
		"step_name": e.sName,
	}).Delete(&models.SStepEnv{}).Error
}
