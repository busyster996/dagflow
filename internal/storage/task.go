package storage

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/busyster996/dagflow/internal/storage/models"
)

type sTask struct {
	*gorm.DB
	tName string

	env IEnv
}

func (t *sTask) Name() string {
	return t.tName
}

func (t *sTask) ClearAll() error {
	if err := t.Remove(); err != nil {
		return err
	}
	if err := t.Env().RemoveAll(); err != nil {
		return err
	}
	list := t.StepList(All)
	for _, v := range list {
		if err := t.Step(v.Name).ClearAll(); err != nil {
			return err
		}
	}
	// 清理build表
	t.Where("task_name", t.tName).Delete(&models.SPipelineBuild{})
	return nil
}

func (t *sTask) Remove() (err error) {
	return t.Where(map[string]interface{}{
		"name": t.tName,
	}).Delete(&models.STask{}).Error
}

func (t *sTask) Kind() (res string, err error) {
	err = t.Model(&models.STask{}).
		Select("kind").
		Where(map[string]interface{}{
			"name": t.tName,
		}).
		Scan(&res).
		Error
	return
}

func (t *sTask) State() (state models.State, err error) {
	err = t.Model(&models.STask{}).
		Select("state").
		Where(map[string]interface{}{
			"name": t.tName,
		}).
		Scan(&state).
		Error
	return
}

func (t *sTask) IsDisable() (disable bool) {
	if t.Model(&models.STask{}).
		Select("disable").
		Where(map[string]interface{}{
			"name": t.tName,
		}).
		Scan(&disable).
		Error != nil {
		return
	}
	return
}

func (t *sTask) Env() IEnv {
	if t.env == nil {
		t.env = &sTaskEnv{
			DB:    t.DB,
			tName: t.tName,
		}
	}
	return t.env
}

func (t *sTask) Timeout() (res time.Duration, err error) {
	err = t.Model(&models.STask{}).
		Select("timeout").
		Where(map[string]interface{}{
			"name": t.tName,
		}).
		Scan(&res).
		Error
	return
}

func (t *sTask) Get() (res *models.STask, err error) {
	res = new(models.STask)
	err = t.Model(&models.STask{}).
		Where(map[string]interface{}{
			"name": t.tName,
		}).
		First(res).
		Error
	return
}

func (t *sTask) UpdateNode(node string) error {
	if node == "" {
		return errors.New("node name can not be empty")
	}
	return t.Model(&models.STask{}).
		Where(map[string]interface{}{
			"name": t.tName,
		}).
		Update("node", node).
		Error
}

func (t *sTask) Update(value *models.STaskUpdate) (err error) {
	if value == nil {
		return
	}
	return t.Model(&models.STask{}).
		Where(map[string]interface{}{
			"name": t.tName,
		}).
		Updates(value).
		Error
}

func (t *sTask) Step(name string) IStep {
	return &sStep{
		DB:    t.DB,
		genv:  t.Env(),
		tName: t.tName,
		sName: name,
	}
}

func (t *sTask) StepCreate(step *models.SStep) (err error) {
	step.TaskName = t.tName
	err = t.Create(step).Error
	if err != nil {
		return fmt.Errorf("save step %d error: %s", step.ID, err)
	}
	return
}

func (t *sTask) StepCount() (res int64) {
	t.Model(&models.SStep{}).Count(&res)
	return
}

func (t *sTask) StepNameList(str string) (res []string) {
	query := t.Model(&models.SStep{}).
		Select("name").
		Order("seq_no ASC").
		Where(map[string]interface{}{
			"task_name": t.tName,
		})
	if str != "" {
		query.Where("name LIKE ?", str)
	}
	query.Find(&res)
	return
}

func (t *sTask) StepStateList(str string) (res map[string]models.State) {
	var steps models.SSteps
	query := t.Model(&models.SStep{}).
		Select("name, state").
		Where(map[string]interface{}{
			"task_name": t.tName,
		})
	if str != "" {
		query.Where("name LIKE ?", str)
	}
	query.Order("seq_no ASC").Find(&steps)
	res = make(map[string]models.State, len(steps))
	for _, v := range steps {
		res[v.Name] = *v.State
	}
	return
}

func (t *sTask) StepList(str string) (res models.SSteps) {
	query := t.Model(&models.SStep{}).
		Where(map[string]interface{}{
			"task_name": t.tName,
		})
	if str != "" {
		query.Where("name LIKE ?", str)
	}
	query.Order("seq_no ASC").Find(&res)
	return
}
