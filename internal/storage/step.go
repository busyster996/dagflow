package storage

import (
	"time"

	"gorm.io/gorm"

	"github.com/busyster996/dagflow/internal/storage/models"
)

type sStep struct {
	*gorm.DB
	genv  IEnv
	tName string
	sName string

	env    IEnv
	depend IDepend
	log    ILog
}

func (s *sStep) Name() string {
	return s.sName
}

func (s *sStep) ClearAll() error {
	if err := s.Remove(); err != nil {
		return err
	}
	if err := s.Env().RemoveAll(); err != nil {
		return err
	}
	if err := s.Depend().RemoveAll(); err != nil {
		return err
	}
	return s.Log().RemoveAll()
}

func (s *sStep) Remove() (err error) {
	return s.Where(map[string]interface{}{
		"task_name": s.tName,
		"name":      s.sName,
	}).Delete(&models.SStep{}).Error
}

func (s *sStep) State() (state models.State, err error) {
	err = s.Model(&models.SStep{}).
		Select("state").
		Where(map[string]interface{}{
			"task_name": s.tName,
			"name":      s.sName,
		}).
		Scan(&state).
		Error
	return
}

func (s *sStep) IsDisable() (disable bool) {
	if s.Model(&models.SStep{}).
		Select("disable").
		Where(map[string]interface{}{
			"task_name": s.tName,
			"name":      s.sName,
		}).
		Scan(&disable).
		Error != nil {
		return
	}
	return
}

func (s *sStep) Env() IEnv {
	if s.env == nil {
		s.env = &sStepEnv{
			DB:    s.DB,
			tName: s.tName,
			sName: s.sName,
		}
	}
	return s.env
}

func (s *sStep) TaskName() string {
	return s.tName
}

func (s *sStep) Timeout() (res time.Duration, err error) {
	err = s.Model(&models.SStep{}).
		Select("timeout").
		Where(map[string]interface{}{
			"task_name": s.tName,
			"name":      s.sName,
		}).
		Scan(&res).
		Error
	return
}

func (s *sStep) Type() (res string, err error) {
	err = s.Model(&models.SStep{}).
		Select("type").
		Where(map[string]interface{}{
			"task_name": s.tName,
			"name":      s.sName,
		}).
		Scan(&res).
		Error
	return
}

func (s *sStep) Content() (res string, err error) {
	err = s.Model(&models.SStep{}).
		Select("content").
		Where(map[string]interface{}{
			"task_name": s.tName,
			"name":      s.sName,
		}).
		Scan(&res).
		Error
	return
}

func (s *sStep) Action() (res string, err error) {
	err = s.Model(&models.SStep{}).
		Select("action").
		Where(map[string]interface{}{
			"task_name": s.tName,
			"name":      s.sName,
		}).
		Scan(&res).
		Error
	return
}

func (s *sStep) Rule() (res string, err error) {
	err = s.Model(&models.SStep{}).
		Select("rule").
		Where(map[string]interface{}{
			"task_name": s.tName,
			"name":      s.sName,
		}).
		Scan(&res).
		Error
	return
}

func (s *sStep) RetryPolicy() (res models.SRetryPolicy, err error) {
	var tmp models.SStep
	err = s.Model(&models.SStep{}).
		Select("retry_policy").
		Where(map[string]interface{}{
			"task_name": s.tName,
			"name":      s.sName,
		}).
		Scan(&tmp).
		Error
	if err != nil {
		return
	}

	res = tmp.RetryPolicy.Data()
	return
}

func (s *sStep) Get() (res *models.SStep, err error) {
	res = new(models.SStep)
	err = s.Model(&models.SStep{}).
		Where(map[string]interface{}{
			"task_name": s.tName,
			"name":      s.sName,
		}).
		First(res).
		Error
	return
}

func (s *sStep) Update(value *models.SStepUpdate) (err error) {
	if value == nil {
		return
	}
	return s.Model(&models.SStep{}).
		Where(map[string]interface{}{
			"task_name": s.tName,
			"name":      s.sName,
		}).
		Updates(value).
		Error
}

func (s *sStep) GlobalEnv() IEnv {
	if s.genv == nil {
		s.genv = &sTaskEnv{
			DB:    s.DB,
			tName: s.tName,
		}
	}
	return s.genv
}

func (s *sStep) Depend() IDepend {
	if s.depend == nil {
		s.depend = &sStepDepend{
			DB:    s.DB,
			tName: s.tName,
			sName: s.sName,
		}
	}
	return s.depend
}

func (s *sStep) Log() ILog {
	if s.log == nil {
		s.log = &sStepLog{
			DB:    s.DB,
			tName: s.tName,
			sName: s.sName,
		}
	}
	return s.log
}
