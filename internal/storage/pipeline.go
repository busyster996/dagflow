package storage

import (
	"gorm.io/gorm"

	"github.com/busyster996/dagflow/internal/storage/models"
)

type sPipeline struct {
	*gorm.DB
	name string
}

func (p *sPipeline) Name() string {
	return p.name
}

func (p *sPipeline) ClearAll() error {
	// clear pipeline
	if err := p.Remove(); err != nil {
		return err
	}
	// clear task
	list, _ := p.Build().List(-1, -1)
	for _, task := range list {
		_ = p.Task(task.TaskName).ClearAll()
	}
	// clear build
	return p.Build().ClearAll()
}

func (p *sPipeline) Remove() (err error) {
	return p.Where(map[string]interface{}{
		"name": p.name,
	}).Delete(&models.SPipeline{}).Error
}

func (p *sPipeline) Build() IPipelineBuild {
	return &sPipelineBuild{
		pName: p.name,
		DB:    p.DB,
	}
}

func (p *sPipeline) Task(name string) ITask {
	return &sTask{
		DB:    p.DB,
		tName: name,
	}
}

func (p *sPipeline) Update(value *models.SPipelineUpdate) (err error) {
	if value == nil {
		return
	}
	return p.Model(&models.SPipeline{}).
		Where(map[string]interface{}{
			"name": p.name,
		}).
		Updates(value).
		Error
}

func (p *sPipeline) Get() (res *models.SPipeline, err error) {
	res = new(models.SPipeline)
	err = p.Model(&models.SPipeline{}).
		Where(map[string]interface{}{
			"name": p.name,
		}).First(res).
		Error
	return
}

func (p *sPipeline) IsDisable() (disable bool) {
	if p.Model(&models.SPipeline{}).
		Select("disable").
		Where(map[string]interface{}{
			"name": p.name,
		}).
		Scan(&disable).
		Error != nil {
		return
	}
	return
}

func (p *sPipeline) TplType() (res string, err error) {
	err = p.Model(&models.SPipeline{}).
		Select("tpl_type").
		Where(map[string]interface{}{
			"name": p.name,
		}).
		Scan(&res).
		Error
	return
}

func (p *sPipeline) Content() (res string, err error) {
	err = p.Model(&models.SPipeline{}).
		Select("content").
		Where(map[string]interface{}{
			"name": p.name,
		}).
		Scan(&res).
		Error
	return
}
